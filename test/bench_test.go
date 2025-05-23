package test

import (
	"crypto/rand"
	"errors"
	"os"
	"testing"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/dgraph-io/badger/v4"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"go.etcd.io/bbolt"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/cursor"
)

const (
	keySize = 511 // voidDB allows 512-byte keys, but LMDB does not
	valSize = 1024

	mapSize = 1 << 40
)

var (
	key [][]byte
	val [][]byte
)

func populateKeyVal(n int) (e error) {
	var (
		i int
	)

	if len(key) >= n {
		return
	}

	key = make([][]byte, n)
	val = make([][]byte, n)

	for i = 0; i < n; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			return
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			return
		}
	}

	return
}

func BenchmarkPopulateKeyVal(b *testing.B) {
	var (
		e error
	)

	e = populateKeyVal(b.N)
	if e != nil {
		b.Fatal(e)
	}

	return
}

func BenchmarkVoidPut(b *testing.B) {
	var (
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				err = txn.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.Update(true, put)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidPutInKeyspace(b *testing.B) {
	const (
		keyspace = "random"
	)

	var (
		cur  *cursor.Cursor
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor([]byte(keyspace))
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = cur.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.Update(true, put)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGet(b *testing.B) {
	var (
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				err = txn.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *voidDB.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				_, err = txn.Get(key[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = void.Update(false, put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGetInKeyspace(b *testing.B) {
	const (
		keyspace = "random"
	)

	var (
		cur  *cursor.Cursor
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor([]byte(keyspace))
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = cur.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor([]byte(keyspace))
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, err = cur.Get(key[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = void.Update(false, put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGetNext(b *testing.B) {
	var (
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				err = txn.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *voidDB.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				_, _, err = txn.GetNext()
				if err != nil {
					return
				}
			}

			return
		}
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = void.Update(false, put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGetNextInKeyspace(b *testing.B) {
	const (
		keyspace = "random"
	)

	var (
		cur  *cursor.Cursor
		e    error
		i    int
		path string
		void *voidDB.Void

		put = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor([]byte(keyspace))
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = cur.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *voidDB.Txn) (err error) {
			cur, err = txn.OpenCursor([]byte(keyspace))
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
	)

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = void.Update(false, put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = void.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBPut(b *testing.B) {
	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenRoot(0)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, 0, 0644)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBPutInDB(b *testing.B) {
	const (
		dbName = "random"
	)

	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMaxDBs(1)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, 0, 0644)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBGet(b *testing.B) {
	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenRoot(0)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *lmdb.Txn) (err error) {
			txn.RawRead = true

			dbi, err = txn.OpenRoot(0)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, err = txn.Get(dbi, key[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, lmdb.NoSync, 0644)
	if e != nil {
		b.Fatal(e)
	}

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBGetInDB(b *testing.B) {
	const (
		dbName = "random"
	)

	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *lmdb.Txn) (err error) {
			txn.RawRead = true

			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_, err = txn.Get(dbi, key[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMaxDBs(1)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, lmdb.NoSync, 0644)
	if e != nil {
		b.Fatal(e)
	}

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBGetNext(b *testing.B) {
	var (
		cur *lmdb.Cursor
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenRoot(0)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *lmdb.Txn) (err error) {
			txn.RawRead = true

			dbi, err = txn.OpenRoot(0)
			if err != nil {
				return
			}

			cur, e = txn.OpenCursor(dbi)
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
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, lmdb.NoSync, 0644)
	if e != nil {
		b.Fatal(e)
	}

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBGetNextInDB(b *testing.B) {
	const (
		dbName = "random"
	)

	var (
		cur *lmdb.Cursor
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		put = func(txn *lmdb.Txn) (err error) {
			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = txn.Put(dbi, key[i], val[i], 0)
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *lmdb.Txn) (err error) {
			txn.RawRead = true

			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			cur, e = txn.OpenCursor(dbi)
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
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMapSize(mapSize)
	if e != nil {
		b.Fatal(e)
	}

	e = env.SetMaxDBs(1)
	if e != nil {
		b.Fatal(e)
	}

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	e = env.Open(tmp, lmdb.NoSync, 0644)
	if e != nil {
		b.Fatal(e)
	}

	e = env.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = env.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkBoltPut(b *testing.B) {
	const (
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		e   error
		i   int
		tmp string

		put = func(txn *bbolt.Tx) (err error) {
			bkt, err = txn.CreateBucket(
				[]byte(bktName),
			)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = bkt.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bdb, e = bbolt.Open(tmp+"/bolt", 0644, nil)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = bdb.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkBoltGet(b *testing.B) {
	const (
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		e   error
		i   int
		tmp string

		put = func(txn *bbolt.Tx) (err error) {
			bkt, err = txn.CreateBucket(
				[]byte(bktName),
			)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = bkt.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *bbolt.Tx) (err error) {
			bkt = txn.Bucket(
				[]byte(bktName),
			)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				_ = bkt.Get(key[i])
			}

			return
		}
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bdb, e = bbolt.Open(tmp+"/bolt", 0644, nil)
	if e != nil {
		b.Fatal(e)
	}

	bdb.NoSync = true

	bdb.NoFreelistSync = true

	e = bdb.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = bdb.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkBoltGetNext(b *testing.B) {
	const (
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		cur *bbolt.Cursor
		e   error
		i   int
		tmp string

		put = func(txn *bbolt.Tx) (err error) {
			bkt, err = txn.CreateBucket(
				[]byte(bktName),
			)
			if err != nil {
				return
			}

			for i = 0; i < b.N; i++ {
				err = bkt.Put(key[i], val[i])
				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *bbolt.Tx) (err error) {
			bkt = txn.Bucket(
				[]byte(bktName),
			)
			if err != nil {
				return
			}

			cur = bkt.Cursor()

			cur.First()

			for i = 1; i < b.N; i++ {
				cur.Next()
			}

			return
		}
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bdb, e = bbolt.Open(tmp+"/bolt", 0644, nil)
	if e != nil {
		b.Fatal(e)
	}

	bdb.NoSync = true

	bdb.NoFreelistSync = true

	e = bdb.Update(put)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	e = bdb.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLevelPut(b *testing.B) {
	var (
		bat leveldb.Batch
		e   error
		i   int
		ldb *leveldb.DB
		tmp string
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	ldb, e = leveldb.OpenFile(tmp+"/level", nil)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		bat.Put(key[i], val[i])
	}

	e = ldb.Write(&bat,
		&opt.WriteOptions{Sync: true},
	)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkLevelGet(b *testing.B) {
	var (
		bat leveldb.Batch
		e   error
		i   int
		ldb *leveldb.DB
		tmp string
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	ldb, e = leveldb.OpenFile(tmp+"/level", nil)
	if e != nil {
		b.Fatal(e)
	}

	for i = 0; i < b.N; i++ {
		bat.Put(key[i], val[i])
	}

	e = ldb.Write(&bat, nil)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		_, e = ldb.Get(key[i], nil)
		if e != nil {
			b.Fatal(e)
		}
	}

	b.StopTimer()

	return
}

func BenchmarkLevelGetNext(b *testing.B) {
	var (
		bat leveldb.Batch
		e   error
		i   int
		itr iterator.Iterator
		ldb *leveldb.DB
		tmp string
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	ldb, e = leveldb.OpenFile(tmp+"/level", nil)
	if e != nil {
		b.Fatal(e)
	}

	for i = 0; i < b.N; i++ {
		bat.Put(key[i], val[i])
	}

	e = ldb.Write(&bat, nil)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	itr = ldb.NewIterator(nil, nil)

	for i = 0; i < b.N; i++ {
		_, _ = itr.Key(), itr.Value()

		e = itr.Error()
		if e != nil {
			b.Fatal(e)
		}

		itr.Next()
	}

	b.StopTimer()

	return
}

func BenchmarkBadgerPut(b *testing.B) {
	var (
		bgr *badger.DB
		e   error
		i   int
		tmp string

		put = func(txn *badger.Txn) (err error) {
			for i = i; i < b.N; i++ {
				err = txn.Set(key[i], val[i])
				if errors.Is(err, badger.ErrTxnTooBig) {
					return nil
				}

				if err != nil {
					return
				}
			}

			return
		}
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bgr, e = badger.Open(
		badger.DefaultOptions(tmp).
			WithSyncWrites(true).
			WithLogger(nil),
	)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	i = 0

	for i < b.N {
		e = bgr.Update(put)
		if e != nil {
			b.Fatal(e)
		}
	}

	b.StopTimer()

	return
}

func BenchmarkBadgerGet(b *testing.B) {
	var (
		bgr *badger.DB
		e   error
		i   int
		tmp string

		put = func(txn *badger.Txn) (err error) {
			for i = i; i < b.N; i++ {
				err = txn.Set(key[i], val[i])
				if errors.Is(err, badger.ErrTxnTooBig) {
					return nil
				}

				if err != nil {
					return
				}
			}

			return
		}

		get = func(txn *badger.Txn) (err error) {
			for i = 0; i < b.N; i++ {
				_, err = txn.Get(key[i])
				if err != nil {
					return
				}
			}

			return
		}
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bgr, e = badger.Open(
		badger.DefaultOptions(tmp).
			WithLogger(nil),
	)
	if e != nil {
		b.Fatal(e)
	}

	i = 0

	for i < b.N {
		e = bgr.Update(put)
		if e != nil {
			b.Fatal(e)
		}
	}

	b.ResetTimer()

	e = bgr.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkBadgerGetNext(b *testing.B) {
	var (
		bgr *badger.DB
		e   error
		i   int
		itr *badger.Iterator
		tmp string

		put = func(txn *badger.Txn) (err error) {
			for i = i; i < b.N; i++ {
				err = txn.Set(key[i], val[i])
				if errors.Is(err, badger.ErrTxnTooBig) {
					return nil
				}

				if err != nil {
					return
				}
			}

			return
		}

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
	)

	tmp, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(tmp)

	bgr, e = badger.Open(
		badger.DefaultOptions(tmp).WithLoggingLevel(badger.ERROR),
	)
	if e != nil {
		b.Fatal(e)
	}

	i = 0

	for i < b.N {
		e = bgr.Update(put)
		if e != nil {
			b.Fatal(e)
		}
	}

	b.ResetTimer()

	e = bgr.View(get)
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkNothing(b *testing.B) {
	var (
		i int
	)

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		continue
	}

	b.StopTimer()

	return
}
