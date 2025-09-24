package database

import (
	"testing"
)

const (
	mapSize = 1 << 40
)

type Constructor func() (Database, error)

type Database interface {
	Get(*testing.B, [][]byte, [][]byte)
	GetNext(*testing.B, [][]byte, [][]byte)
	Put(*testing.B, [][]byte, [][]byte, bool)

	UseKeyspace([]byte)
	UseShortLivedTxns()

	Close() error
}
