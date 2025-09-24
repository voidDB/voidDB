package database

import (
	"os"
	"testing"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/cursor"
)

type void struct {
	path      string
	void      *voidDB.Void
	keyspace  []byte
	shortTxns bool
}

func NewVoidDB() (Database, error) {
	return newVoidDB()
}

func newVoidDB() (db *void, e error) {
	db = new(void)

	db.path, e = os.MkdirTemp("", "")
	if e != nil {
		return
	}

	db.void, e = voidDB.NewVoid(db.path+"/void", mapSize)
	if e != nil {
		return
	}

	return
}

func (db *void) Close() error {
	return os.RemoveAll(db.path)
}

func (db *void) UseKeyspace(keyspace []byte) {
	db.keyspace = keyspace

	return
}

func (db *void) UseShortLivedTxns() {
	db.shortTxns = true

	return
}

func (db *void) Put(b *testing.B, keys, values [][]byte, sync bool) {
	var (
		cur *cursor.Cursor
		e   error
		i   int
		put func(*voidDB.Txn) error
	)

	switch db.shortTxns {
	case true:
		put = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor(db.keyspace)
			if err != nil {
				return
			}

			return txn.Put(keys[i], values[i])
		}

		for i = 0; i < b.N; i++ {
			e = db.void.Update(sync, put)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		put = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor(db.keyspace)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = cur.Put(keys[i], values[i])
				if err != nil {
					return
				}
			}

			return
		}

		e = db.void.Update(sync, put)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *void) Get(b *testing.B, keys, values [][]byte) {
	var (
		cur *cursor.Cursor
		e   error
		get func(*voidDB.Txn) error
		i   int
	)

	switch db.shortTxns {
	case true:
		get = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor(db.keyspace)
			if err != nil {
				return
			}

			_, err = txn.Get(keys[i])
			if err != nil {
				return
			}

			return
		}

		for i = 0; i < b.N; i++ {
			e = db.void.View(get)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		get = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor(db.keyspace)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, err = cur.Get(keys[i])
				if err != nil {
					return
				}
			}

			return
		}

		e = db.void.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *void) GetNext(b *testing.B, keys, values [][]byte) {
	var (
		cur *cursor.Cursor
		e   error
		get func(*voidDB.Txn) error
		i   int
	)

	switch db.shortTxns {
	case true:
		b.Fatal("GetNext test not implemented for short-lived transactions")

	default:
		get = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor(db.keyspace)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, _, err = cur.GetNext()
				if err != nil {
					return
				}
			}

			return
		}

		e = db.void.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}
