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

func NewLoop(ctx context.Context, link *Link, bucketName, uploadName string,
	uplinkPeriod, downlinkPeriod time.Duration,
) (
	loop *Loop,
) {
	loop = &Loop{
		context:        ctx,
		link:           link,
		bucketName:     bucketName,
		uploadName:     uploadName,
		uplinkTicker:   time.NewTicker(uplinkPeriod),
		downlinkTicker: time.NewTicker(downlinkPeriod),
	}

	return
}

func (loop *Loop) Uplink() (e error) {
	var (
		txnID     int
		txnIDLast int

		getTxnID = func(txn *voidDB.Txn) error {
			txnID = txn.SerialNumber()

			return nil
		}
	)

	for {
		select {
		case <-loop.context.Done():
			return loop.context.Err()

		case <-loop.uplinkTicker.C:
			e = loop.link.void.View(getTxnID)
			if e != nil {
				log.Error().Err(e)

				continue
			}

			if txnID <= txnIDLast+1 {
				continue
			}

			e = loop.link.Uplink(loop.context, loop.bucketName, loop.uploadName)
			if e != nil {
				log.Error().Err(e)

				continue
			}

			txnIDLast = txnID
		}
	}
}

func (loop *Loop) Downlink() (e error) {
	var (
		objectChan <-chan minio.ObjectInfo
		objectInfo minio.ObjectInfo

		objectETag = make(map[string]string)
	)

	for {
		select {
		case <-loop.context.Done():
			return loop.context.Err()

		case <-loop.downlinkTicker.C:
			objectChan = loop.link.client.ListObjects(loop.context,
				loop.bucketName, minio.ListObjectsOptions{},
			)

			for objectInfo = range objectChan {
				e = loop.context.Err()
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

				e = loop.link.Downlink(loop.context, loop.bucketName,
					objectInfo.Key,
				)
				if e != nil {
					log.Error().Err(e)

					continue
				}

				objectETag[objectInfo.Key] = objectInfo.ETag
			}
		}
	}
}
