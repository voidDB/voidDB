package database

import (
	"fmt"
	"os"
	"testing"

	"go.etcd.io/bbolt"
)

type bolt struct {
	path      string
	db        *bbolt.DB
	keyspace  []byte
	shortTxns bool
}

func NewBoltDB() (Database, error) {
	return newBoltDB()
}

func newBoltDB() (db *bolt, e error) {
	db = new(bolt)

	db.path, e = os.MkdirTemp("", "")
	if e != nil {
		return
	}

	db.db, e = bbolt.Open(db.path+"/bolt", 0644, nil)
	if e != nil {
		return
	}

	return
}

func (db *bolt) Close() error {
	return os.RemoveAll(db.path)
}

func (db *bolt) UseKeyspace(keyspace []byte) {
	db.keyspace = keyspace

	return
}

func (db *bolt) UseShortLivedTxns() {
	db.shortTxns = true

	return
}

func (db *bolt) Put(b *testing.B, keys, values [][]byte, sync bool) {
	var (
		bkt *bbolt.Bucket
		e   error
		i   int
		put func(*bbolt.Tx) error
	)

	switch sync {
	case true:
		db.db.NoSync = false

		db.db.NoFreelistSync = false

	case false:
		db.db.NoSync = true

		db.db.NoFreelistSync = true
	}

	switch db.shortTxns {
	case true:
		put = func(txn *bbolt.Tx) (err error) {
			bkt, err = txn.CreateBucketIfNotExists(db.keyspace)
			if err != nil {
				return
			}

			return bkt.Put(keys[i], values[i])
		}

		for i = 0; i < b.N; i++ {
			e = db.db.Update(put)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		put = func(txn *bbolt.Tx) (err error) {
			bkt, err = txn.CreateBucketIfNotExists(db.keyspace)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = bkt.Put(keys[i], values[i])
				if err != nil {
					return
				}
			}

			return
		}

		e = db.db.Update(put)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *bolt) Get(b *testing.B, keys, values [][]byte) {
	var (
		bkt *bbolt.Bucket
		e   error
		get func(*bbolt.Tx) error
		i   int
	)

	switch db.shortTxns {
	case true:
		get = func(txn *bbolt.Tx) (err error) {
			bkt = txn.Bucket(db.keyspace)
			if bkt == nil {
				err = fmt.Errorf("Bucket does not exist")
			}

			if bkt.Get(keys[i]) == nil {
				err = fmt.Errorf("Key not found")
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
		get = func(txn *bbolt.Tx) (err error) {
			bkt = txn.Bucket(db.keyspace)
			if bkt == nil {
				err = fmt.Errorf("Bucket does not exist")
			}

			for i = 0; i < b.N; i++ {
				if bkt.Get(keys[i]) == nil {
					err = fmt.Errorf("Key not found")
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

func (db *bolt) GetNext(b *testing.B, keys, values [][]byte) {
	var (
		bkt *bbolt.Bucket
		cur *bbolt.Cursor
		e   error
		get func(*bbolt.Tx) error
		i   int
	)

	switch db.shortTxns {
	case true:
		b.Fatal("GetNext test not implemented for short-lived transactions")

	default:
		get = func(txn *bbolt.Tx) (err error) {
			bkt = txn.Bucket(db.keyspace)
			if bkt == nil {
				err = fmt.Errorf("Bucket does not exist")
			}

			cur = bkt.Cursor()

			cur.First()

			for i = 1; i < b.N; i++ {
				cur.Next()
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
