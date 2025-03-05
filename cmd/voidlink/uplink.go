package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"os/exec"

	"github.com/minio/minio-go"
	"golang.org/x/sync/errgroup"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
	"github.com/voidDB/voidDB/link"
)

func (l *Link) Uplink(ctx context.Context, bucketName, objectName string) (
	e error,
) {
	var (
		group *errgroup.Group
	)

	l.mutex.Lock()

	defer l.renew()

	group, l.context = errgroup.WithContext(ctx)

	l.bucketName = bucketName
	l.objectName = objectName

	l.txn, e = l.void.BeginTxn(true, false)
	if e != nil {
		return
	}

	defer l.txn.Abort()

	group.Go(l.pick)
	group.Go(l.serialise)
	group.Go(l.compress)
	group.Go(l.upload)

	return group.Wait()
}

func (l *Link) pick() (e error) {
	var (
		iterRecords = func(proto *cursor.Cursor) (err error) {
			var (
				cur Cursor = proto.ToLinkCursor()

				key      []byte
				metadata link.Metadata
				value    []byte

				record Record
			)

			for {
				err = l.context.Err()
				if err != nil {
					return
				}

				key, value, metadata, err = cur.GetNext(0)

				if err != nil && !errors.Is(err, common.ErrorDeleted) {
					break
				}

				record = Record{
					key:      key,
					metadata: metadata.Timestamp(),
					value:    value,
				}

				l.channel <- record
			}

			return
		}

		iterKeyspaces = func(txn *voidDB.Txn) (err error) {
			var (
				cur      *cursor.Cursor
				keyspace []byte
			)

			for {
				err = l.context.Err()
				if err != nil {
					return
				}

				keyspace, _, err = txn.GetNext()
				if err != nil {
					break
				}

				cur, err = txn.OpenCursor(keyspace)
				if err != nil {
					return
				}

				l.channel <- Record{key: keyspace}

				err = iterRecords(cur)
				if !errors.Is(err, common.ErrorNotFound) {
					break
				}
			}

			if errors.Is(err, common.ErrorNotFound) {
				err = nil
			}

			return
		}
	)

	e = iterKeyspaces(l.txn)
	if e != nil {
		return fmt.Errorf("error occurred while picking: %w", e)
	}

	close(l.channel)

	return
}

func (l *Link) serialise() (e error) {
	var (
		encoder *Encoder = NewEncoder(l.writer0,
			fnv.New32a(),
		)

		record Record
		recvOK bool
	)

	_, e = l.writer0.Write(
		[]byte(linkMagic),
	)
	if e != nil {
		return
	}

	e = binary.Write(l.writer0, binary.BigEndian,
		uint64(version),
	)
	if e != nil {
		return
	}

	for {
		select {
		case <-l.context.Done():
			return l.context.Err()

		case record, recvOK = <-l.channel:
			break
		}

		e = encoder.Encode(record.key, record.metadata, record.value)
		if e != nil {
			return fmt.Errorf("error occurred during serialisation: %w", e)
		}

		if !recvOK {
			break
		}
	}

	l.writer0.Close()

	return
}

func (l *Link) compress() (e error) {
	var (
		command *exec.Cmd = exec.CommandContext(l.context,
			"zstd", "--stdout", "-",
		)
	)

	command.Stdin, command.Stdout = l.reader0, l.writer1

	e = command.Run()
	if e != nil {
		return fmt.Errorf("error occurred during compression: %s",
			e.(*exec.ExitError).Stderr,
		)
	}

	l.writer1.Close()

	return
}

func (l *Link) upload() (e error) {
	const (
		objectSize = -1 // do multipart Put until EOF
	)

	_, e = l.client.PutObjectWithContext(l.context, l.bucketName, l.objectName,
		l.reader1, objectSize, minio.PutObjectOptions{},
	)
	if e != nil {
		return fmt.Errorf("error occurred during upload: %w", e)
	}

	return
}
