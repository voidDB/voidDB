package database

import (
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type level struct {
	path      string
	db        *leveldb.DB
	keyspace  string
	shortTxns bool
}

func NewLevelDB() (Database, error) {
	return newLevelDB()
}

func newLevelDB() (db *level, e error) {
	db = new(level)

	db.path, e = os.MkdirTemp("", "")
	if e != nil {
		return
	}

	db.db, e = leveldb.OpenFile(db.path+"/level", nil)
	if e != nil {
		return
	}

	return
}

func (db *level) Close() error {
	return os.RemoveAll(db.path)
}

func (db *level) UseKeyspace(keyspace []byte) {
	return
}

func (db *level) UseShortLivedTxns() {
	db.shortTxns = true

	return
}

func (db *level) Put(b *testing.B, keys, values [][]byte, sync bool) {
	var (
		bat leveldb.Batch
		e   error
		i   int
	)

	switch db.shortTxns {
	case true:
		for i = 0; i < b.N; i++ {
			e = db.db.Put(keys[i], values[i],
				&opt.WriteOptions{Sync: sync},
			)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		for i = 0; i < b.N; i++ {
			bat.Put(keys[i], values[i])
		}

		e = db.db.Write(&bat,
			&opt.WriteOptions{Sync: sync},
		)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *level) Get(b *testing.B, keys, values [][]byte) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		_, e = db.db.Get(keys[i], nil)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *level) GetNext(b *testing.B, keys, values [][]byte) {
	var (
		e   error
		i   int
		itr iterator.Iterator
	)

	switch db.shortTxns {
	case true:
		b.Fatal("GetNext test not implemented for short-lived transactions")

	default:
		itr = db.db.NewIterator(nil, nil)

		for i = 0; i < b.N; i++ {
			_, _ = itr.Key(), itr.Value()

			e = itr.Error()
			if e != nil {
				b.Fatal(e)
			}

			itr.Next()
		}
	}

	return
}
