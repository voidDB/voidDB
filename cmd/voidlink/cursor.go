package main

import (
	"github.com/voidDB/voidDB/link"
)

type Cursor interface {
	GetNext(int) ([]byte, []byte, link.Metadata, error)
	Get([]byte) (link.Metadata, error)
	Put([]byte, []byte, link.Metadata) error
	Del(link.Metadata) error
}
