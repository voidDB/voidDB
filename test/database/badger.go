package database

import (
	"errors"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v4"
)

type badgerDB struct {
	path      string
	db        *badger.DB
	keyspace  string
	shortTxns bool
}

func NewBadgerDB() (Database, error) {
	return newBadgerDB()
}

func newBadgerDB() (db *badgerDB, e error) {
	db = new(badgerDB)

	db.path, e = os.MkdirTemp("", "")
	if e != nil {
		return
	}

	db.db, e = badger.Open(
		badger.DefaultOptions(db.path).
			WithSyncWrites(true).
			WithLogger(nil),
	)
	if e != nil {
		return
	}

	return
}

func (db *badgerDB) Close() error {
	return os.RemoveAll(db.path)
}

func (db *badgerDB) UseKeyspace(keyspace []byte) {
	return
}

func (db *badgerDB) UseShortLivedTxns() {
	db.shortTxns = true

	return
}

func (db *badgerDB) Put(b *testing.B, keys, values [][]byte, sync bool) {
	var (
		e   error
		i   int
		put func(*badger.Txn) error
	)

	switch db.shortTxns {
	case true:
		put = func(txn *badger.Txn) (err error) {
			return txn.Set(keys[i], values[i])
		}

		for i = 0; i < b.N; i++ {
			e = db.db.Update(put)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		put = func(txn *badger.Txn) (err error) {
			for i = i; i < b.N; i++ {
				err = txn.Set(keys[i], values[i])
				if errors.Is(err, badger.ErrTxnTooBig) {
					return nil
				}

				if err != nil {
					return
				}
			}

			return
		}

		for i < b.N {
			e = db.db.Update(put)
			if e != nil {
				b.Fatal(e)
			}
		}
	}

	return
}

func (db *badgerDB) Get(b *testing.B, keys, values [][]byte) {
	var (
		e   error
		get func(*badger.Txn) error
		i   int
	)

	switch db.shortTxns {
	case true:
		get = func(txn *badger.Txn) (err error) {
			_, err = txn.Get(keys[i])
			if err != nil {
				return
			}

			return
		}

		for i = 0; i < b.N; i++ {
			e = db.db.View(get)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		get = func(txn *badger.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				_, err = txn.Get(keys[i])
				if err != nil {
					return
				}
			}

			return
		}

		e = db.db.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *badgerDB) GetNext(b *testing.B, keys, values [][]byte) {
	var (
		e   error
		get func(*badger.Txn) error
		i   int
		itr *badger.Iterator
	)

	switch db.shortTxns {
	case true:
		b.Fatal("GetNext test not implemented for short-lived transactions")

	default:
		get = func(txn *badger.Txn) (err error) {
			itr = txn.NewIterator(badger.DefaultIteratorOptions)

			itr.Rewind()

			for i = 0; i < b.N; i++ {
				itr.Item()

				itr.Next()
			}

			itr.Close()

			return
		}

		e = db.db.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}
