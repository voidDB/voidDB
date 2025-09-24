package database

import (
	"os"
	"testing"

	"github.com/PowerDNS/lmdb-go/lmdb"
)

type lMDB struct {
	path      string
	env       *lmdb.Env
	keyspace  string
	shortTxns bool
}

func NewLMDB() (Database, error) {
	return newLMDB()
}

func newLMDB() (db *lMDB, e error) {
	db = new(lMDB)

	db.path, e = os.MkdirTemp("", "")
	if e != nil {
		return
	}

	db.env, e = lmdb.NewEnv()
	if e != nil {
		return
	}

	e = db.env.SetMapSize(mapSize)
	if e != nil {
		return
	}

	e = db.env.SetMaxDBs(1)
	if e != nil {
		return
	}

	e = db.env.Open(db.path, 0, 0644)
	if e != nil {
		return
	}

	return
}

func (db *lMDB) Close() error {
	return os.RemoveAll(db.path)
}

func (db *lMDB) UseKeyspace(keyspace []byte) {
	db.keyspace = string(keyspace)

	return
}

func (db *lMDB) UseShortLivedTxns() {
	db.shortTxns = true

	return
}

func (db *lMDB) Put(b *testing.B, keys, values [][]byte, sync bool) {
	var (
		cur *lmdb.Cursor
		dbi lmdb.DBI
		e   error
		i   int
		put func(*lmdb.Txn) error
	)

	switch db.shortTxns {
	case true:
		put = func(txn *lmdb.Txn) (err error) {
			switch db.keyspace {
			case "":
				dbi, err = txn.OpenRoot(0)

			default:
				dbi, err = txn.OpenDBI(db.keyspace, lmdb.Create)
			}
			if err != nil {
				return
			}

			cur, err = txn.OpenCursor(dbi)
			if err != nil {
				return
			}

			return cur.Put(keys[i], values[i], 0)
		}

		for i = 0; i < b.N; i++ {
			e = db.env.Update(put)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		put = func(txn *lmdb.Txn) (err error) {
			switch db.keyspace {
			case "":
				dbi, err = txn.OpenRoot(0)

			default:
				dbi, err = txn.OpenDBI(db.keyspace, lmdb.Create)
			}
			if err != nil {
				return
			}

			cur, err = txn.OpenCursor(dbi)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = cur.Put(keys[i], values[i], 0)
				if err != nil {
					return
				}
			}

			return
		}

		e = db.env.Update(put)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *lMDB) Get(b *testing.B, keys, values [][]byte) {
	var (
		cur *lmdb.Cursor
		e   error
		get func(*lmdb.Txn) error
		i   int
		dbi lmdb.DBI
	)

	switch db.shortTxns {
	case true:
		get = func(txn *lmdb.Txn) (err error) {
			switch db.keyspace {
			case "":
				dbi, err = txn.OpenRoot(0)

			default:
				dbi, err = txn.OpenDBI(db.keyspace, 0)
			}
			if err != nil {
				return
			}

			cur, err = txn.OpenCursor(dbi)
			if err != nil {
				return
			}

			_, _, err = cur.Get(keys[i], nil, 0)
			if err != nil {
				return
			}

			return
		}

		for i = 0; i < b.N; i++ {
			e = db.env.View(get)
			if e != nil {
				b.Fatal(e)
			}
		}

	default:
		get = func(txn *lmdb.Txn) (err error) {
			switch db.keyspace {
			case "":
				dbi, err = txn.OpenRoot(0)

			default:
				dbi, err = txn.OpenDBI(db.keyspace, 0)
			}
			if err != nil {
				return
			}

			cur, err = txn.OpenCursor(dbi)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, _, err = cur.Get(keys[i], nil, 0)
				if err != nil {
					return
				}
			}

			return
		}

		e = db.env.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}

func (db *lMDB) GetNext(b *testing.B, keys, values [][]byte) {
	var (
		cur *lmdb.Cursor
		e   error
		get func(*lmdb.Txn) error
		i   int
		dbi lmdb.DBI
	)

	switch db.shortTxns {
	case true:
		b.Fatal("GetNext test not implemented for short-lived transactions")

	default:
		get = func(txn *lmdb.Txn) (err error) {
			switch db.keyspace {
			case "":
				dbi, err = txn.OpenRoot(0)

			default:
				dbi, err = txn.OpenDBI(db.keyspace, 0)
			}
			if err != nil {
				return
			}

			cur, err = txn.OpenCursor(dbi)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, _, err = cur.Get(nil, nil, lmdb.Next)
				if err != nil {
					return
				}
			}

			return
		}

		e = db.env.View(get)
		if e != nil {
			b.Fatal(e)
		}
	}

	return
}
