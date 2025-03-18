package main

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"

	"github.com/voidDB/voidDB"
)

type Loop struct {
	context        context.Context
	link           *Link
	bucketName     string
	uploadName     string
	uplinkTicker   *time.Ticker
	downlinkTicker *time.Ticker
}

func NewLoop(link *Link, bucketName, uploadName string,
	uplinkPeriod, downlinkPeriod time.Duration,
) (
	loop *Loop,
) {
	loop = &Loop{
		link:       link,
		bucketName: bucketName,
		uploadName: uploadName,
	}

	if uplinkPeriod > 0 {
		loop.uplinkTicker = time.NewTicker(uplinkPeriod)
	}

	if downlinkPeriod > 0 {
		loop.downlinkTicker = time.NewTicker(downlinkPeriod)
	}

	return
}

func (loop *Loop) Uplink(ctx context.Context) (e error) {
	var (
		txnID     int
		txnIDLast int

		getTxnID = func(txn *voidDB.Txn) error {
			txnID = txn.SerialNumber() - 1

			return nil
		}
	)

	if loop.uplinkTicker == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-loop.uplinkTicker.C:
			e = loop.link.void.View(getTxnID)
			if e != nil {
				log.Error().Err(e)

				continue
			}

			if txnID <= txnIDLast {
				continue
			}

			e = loop.link.Uplink(ctx, loop.bucketName, loop.uploadName)
			if e != nil {
				log.Error().Err(e)

				continue
			}

			txnIDLast = txnID
		}
	}
}

func (loop *Loop) Downlink(ctx context.Context) (e error) {
	var (
		eTag       string
		objectChan <-chan minio.ObjectInfo
		objectInfo minio.ObjectInfo

		objectETag = make(map[string]string)
	)

	if loop.downlinkTicker == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-loop.downlinkTicker.C:
			objectInfo, _ = loop.link.client.StatObject(ctx,
				loop.bucketName, loop.uploadName, minio.StatObjectOptions{},
			)

			if objectInfo.ETag != "" {
				objectETag[loop.uploadName] = objectInfo.ETag
			}

			objectChan = loop.link.client.ListObjects(ctx,
				loop.bucketName, minio.ListObjectsOptions{},
			)

			for objectInfo = range objectChan {
				e = ctx.Err()
				if e != nil {
					continue // drain objectChan, per Minio documentation
				}

				e = objectInfo.Err

				switch {
				case e != nil:
					log.Error().Err(e)

					continue

				case objectInfo.Key == loop.uploadName:
					continue

				case objectETag[objectInfo.Key] == objectInfo.ETag:
					continue
				}

				for _, eTag = range objectETag {
					if eTag == objectInfo.ETag {
						goto end
					}
				}

				e = loop.link.Downlink(ctx, loop.bucketName, objectInfo.Key)
				if e != nil {
					log.Error().Err(e)

					continue
				}

			end:
				objectETag[objectInfo.Key] = objectInfo.ETag
			}
		}
	}
}
