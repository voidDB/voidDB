package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os/exec"
	"syscall"

	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
	"github.com/voidDB/voidDB/link"
)

func (l *Link) Downlink(ctx context.Context, bucketName, objectName string) (
	e error,
) {
	var (
		group *errgroup.Group
	)

	l.mutex.Lock()

	defer l.mutex.Unlock()

	l.renew()

	group, l.context = errgroup.WithContext(ctx)

	l.bucketName = bucketName
	l.objectName = objectName

	group.Go(l.download)
	group.Go(l.decompress)
	group.Go(l.deserialise)
	group.Go(l.reconcile)

	return group.Wait()
}

func (l *Link) download() (e error) {
	var (
		object *minio.Object
	)

	object, e = l.client.GetObject(l.context, l.bucketName, l.objectName,
		minio.GetObjectOptions{},
	)
	if e != nil {
		goto bad
	}

	defer object.Close()

	_, e = io.Copy(l.writer0, object)
	if e != nil {
		goto bad
	}

	l.writer0.Close()

	return

bad:
	return fmt.Errorf("error occurred during download: %w", e)
}

func (l *Link) decompress() (e error) {
	var (
		command *exec.Cmd = exec.CommandContext(l.context,
			"unzstd", "--stdout", "-",
		)
	)

	command.Stdin, command.Stdout = l.reader0, l.writer1

	e = command.Run()
	if e != nil {
		return fmt.Errorf("error occurred during decompression: %w", e)
	}

	l.writer1.Close()

	return
}

func (l *Link) deserialise() (e error) {
	var (
		decoder *Decoder = NewDecoder(l.reader1,
			fnv.New32a(),
		)

		magic []byte = make([]byte,
			len(linkMagic),
		)

		record Record
		verNo  uint64
	)

	_, e = l.reader1.Read(magic)
	if e != nil {
		return
	}

	if !bytes.Equal(magic, []byte(linkMagic)) {
		return syscall.EPROTO
	}

	e = binary.Read(l.reader1, binary.BigEndian, &verNo)
	if e != nil {
		return
	}

	if verNo != version {
		return syscall.EPROTO
	}

	for {
		e = l.context.Err()
		if e != nil {
			return
		}

		record.key, record.metadata, record.value, e = decoder.Decode()
		if e != nil {
			return fmt.Errorf("error occurred during deserialisation: %w", e)
		}

		if len(record.key)+len(record.metadata)+len(record.value) == 0 {
			break
		}

		l.channel <- record
	}

	close(l.channel)

	return nil
}

func (l *Link) reconcile() (e error) {
	var (
		cur    Cursor
		meta   link.Metadata
		proto  *cursor.Cursor
		record Record
		recvOK bool

		merge = func(txn *voidDB.Txn) (err error) {
			for {
				select {
				case <-l.context.Done():
					return l.context.Err()

				case record, recvOK = <-l.channel:
					break
				}

				if !recvOK {
					break
				}

				switch {
				case len(record.metadata) == 0:
					proto, err = txn.OpenCursor(record.key)
					if err != nil {
						return
					}

					cur = proto.ToLinkCursor()

					continue

				default:
					meta, err = cur.Get(record.key)
				}

				switch {
				case errors.Is(err, common.ErrorNotFound):
					break

				case err != nil:
					return

				case bytes.Compare(record.metadata, meta.Timestamp()) > 0:
					break

				default:
					continue
				}

				meta = link.NewMetadata(record.metadata, 0)

				switch {
				case len(record.value) > 0:
					err = cur.Put(record.key, record.value, meta)

				default:
					err = cur.Del(meta)
				}
				if err != nil {
					return
				}
			}

			return
		}
	)

	e = l.void.Update(true, merge)
	if e != nil {
		return fmt.Errorf("error occurred during reconciliation: %w", e)
	}

	return
}
