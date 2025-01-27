package test

import (
	"crypto/rand"
	"errors"
	"os"
	"testing"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/dgraph-io/badger/v4"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"go.etcd.io/bbolt"

	"github.com/voidDB/voidDB"
)

func BenchmarkVoidPut(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 30 // 1 GiB
	)

	var (
		e    error
		i    int
		path string
		txnW *voidDB.Txn
		void *voidDB.Void

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	txnW, e = void.BeginTxn(false, true)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		e = txnW.Put(key[i], val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	e = txnW.Commit()
	if e != nil {
		b.Fatal(e)
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGet(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 30 // 1 GiB
	)

	var (
		e    error
		i    int
		path string
		txnR *voidDB.Txn
		txnW *voidDB.Txn
		void *voidDB.Void

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	txnW, e = void.BeginTxn(false, false)
	if e != nil {
		b.Fatal(e)
	}

	for i = 0; i < b.N; i++ {
		e = txnW.Put(key[i], val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	e = txnW.Commit()
	if e != nil {
		b.Fatal(e)
	}

	txnR, e = void.BeginTxn(true, false)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		_, e = txnR.Get(key[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	b.StopTimer()

	return
}

func BenchmarkVoidGetNext(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 30 // 1 GiB
	)

	var (
		e    error
		i    int
		path string
		txnR *voidDB.Txn
		txnW *voidDB.Txn
		void *voidDB.Void

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	path, e = os.MkdirTemp("", "")
	if e != nil {
		b.Fatal(e)
	}

	defer os.RemoveAll(path)

	void, e = voidDB.NewVoid(path+"/void", mapSize)
	if e != nil {
		b.Fatal(e)
	}

	txnW, e = void.BeginTxn(false, false)
	if e != nil {
		b.Fatal(e)
	}

	for i = 0; i < b.N; i++ {
		e = txnW.Put(key[i], val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

	e = txnW.Commit()
	if e != nil {
		b.Fatal(e)
	}

	txnR, e = void.BeginTxn(true, false)
	if e != nil {
		b.Fatal(e)
	}

	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		_, _, e = txnR.GetNext()
		if e != nil {
			b.Fatal(e)
		}
	}

	b.StopTimer()

	return
}

func BenchmarkLMDBPut(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 31 // 2 GiB
	)

	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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

func BenchmarkLMDBGet(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 31 // 2 GiB
	)

	var (
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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

func BenchmarkLMDBGetNext(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		mapSize = 1 << 31 // 2 GiB
	)

	var (
		cur *lmdb.Cursor
		dbi lmdb.DBI
		e   error
		env *lmdb.Env
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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

func BenchmarkBoltPut(b *testing.B) {
	const (
		keySize = 511
		valSize = 4096
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		e   error
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
		keySize = 511
		valSize = 4096
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		e   error
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
		keySize = 511
		valSize = 4096
		bktName = "random"
	)

	var (
		bdb *bbolt.DB
		bkt *bbolt.Bucket
		cur *bbolt.Cursor
		e   error
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bat leveldb.Batch
		e   error
		i   int
		ldb *leveldb.DB
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bat leveldb.Batch
		e   error
		i   int
		ldb *leveldb.DB
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bat leveldb.Batch
		e   error
		i   int
		itr iterator.Iterator
		ldb *leveldb.DB
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)
	)

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bgr *badger.DB
		e   error
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bgr *badger.DB
		e   error
		i   int
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
	const (
		keySize = 511
		valSize = 4096
	)

	var (
		bgr *badger.DB
		e   error
		i   int
		itr *badger.Iterator
		tmp string

		key = make([][]byte, b.N)
		val = make([][]byte, b.N)

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

	for i = 0; i < b.N; i++ {
		key[i] = make([]byte, keySize)

		_, e = rand.Read(key[i])
		if e != nil {
			b.Fatal(e)
		}

		val[i] = make([]byte, valSize)

		_, e = rand.Read(val[i])
		if e != nil {
			b.Fatal(e)
		}
	}

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
