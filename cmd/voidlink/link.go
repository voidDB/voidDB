package main

import (
	"context"
	"io"
	"sync"

	"github.com/minio/minio-go"

	"github.com/voidDB/voidDB"
)

const (
	linkMagic = "voidLINK"
	version   = 0
)

type Link struct {
	void   *voidDB.Void
	client *minio.Client

	mutex      sync.Mutex
	context    context.Context
	bucketName string
	objectName string
	txn        *voidDB.Txn

	channel chan Record
	writer0 io.WriteCloser
	reader0 io.ReadCloser
	writer1 io.WriteCloser
	reader1 io.ReadCloser
}

func NewLink(void *voidDB.Void, client *minio.Client) (l *Link) {
	l = &Link{
		void:    void,
		client:  client,
		channel: make(chan Record),
	}

	l.reader0, l.writer0 = io.Pipe()
	l.reader1, l.writer1 = io.Pipe()

	return
}

func (l *Link) renew() {
	*l = *NewLink(l.void, l.client)

	return
}

type Record struct {
	key      []byte
	metadata []byte
	value    []byte
}
