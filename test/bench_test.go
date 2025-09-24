package test

import (
	"crypto/rand"
	"testing"

	"test/database"
)

const (
	keySize = 511 // voidDB allows 512-byte keys, but LMDB does not
	valSize = 1024
)

var (
	keyspace = []byte("random")

	key [][]byte
	val [][]byte
)

func prepareToBenchmark(b *testing.B, constructor database.Constructor,
	keyspace []byte, shortTxns bool, prepopulate bool,
) (
	db database.Database,
) {
	var (
		e error
	)

	e = populateKeyVal(b.N)
	if e != nil {
		b.Fatal(e)
	}

	db, e = constructor()
	if e != nil {
		b.Fatal(e)
	}

	if keyspace != nil {
		db.UseKeyspace(keyspace)
	}

	if shortTxns {
		db.UseShortLivedTxns()
	}

	if prepopulate {
		db.Put(b, key, val, false)
	}

	return
}

func benchmarkPut(b *testing.B, constructor database.Constructor,
	keyspace []byte, shortTxns bool,
) {
	var (
		db database.Database
	)

	db = prepareToBenchmark(b, constructor, keyspace, shortTxns, false)

	defer db.Close()

	b.ResetTimer()

	db.Put(b, key, val, true)

	b.StopTimer()
}

func benchmarkGet(b *testing.B, constructor database.Constructor,
	keyspace []byte, shortTxns bool,
) {
	var (
		db database.Database
	)

	db = prepareToBenchmark(b, constructor, keyspace, shortTxns, true)

	defer db.Close()

	b.ResetTimer()

	db.Get(b, key, val)

	b.StopTimer()
}

func benchmarkGetNext(b *testing.B, constructor database.Constructor,
	keyspace []byte, shortTxns bool,
) {
	var (
		db database.Database
	)

	db = prepareToBenchmark(b, constructor, keyspace, shortTxns, true)

	defer db.Close()

	b.ResetTimer()

	db.GetNext(b, key, val)

	b.StopTimer()
}

func BenchmarkVoidPut(b *testing.B) {
	benchmarkPut(b, database.NewVoidDB, nil, false)

	return
}

func BenchmarkVoidPutInKeyspace(b *testing.B) {
	benchmarkPut(b, database.NewVoidDB, keyspace, false)

	return
}

func BenchmarkVoidGet(b *testing.B) {
	benchmarkGet(b, database.NewVoidDB, nil, false)

	return
}

func BenchmarkVoidGetInKeyspace(b *testing.B) {
	benchmarkGet(b, database.NewVoidDB, keyspace, false)

	return
}

func BenchmarkVoidGetNext(b *testing.B) {
	benchmarkGetNext(b, database.NewVoidDB, nil, false)

	return
}

func BenchmarkVoidGetNextInKeyspace(b *testing.B) {
	benchmarkGetNext(b, database.NewVoidDB, keyspace, false)

	return
}

func BenchmarkLMDBPut(b *testing.B) {
	benchmarkPut(b, database.NewLMDB, nil, false)

	return
}

func BenchmarkLMDBPutInDB(b *testing.B) {
	benchmarkPut(b, database.NewLMDB, keyspace, false)

	return
}

func BenchmarkLMDBGet(b *testing.B) {
	benchmarkGet(b, database.NewLMDB, nil, false)

	return
}

func BenchmarkLMDBGetInDB(b *testing.B) {
	benchmarkGet(b, database.NewLMDB, keyspace, false)

	return
}

func BenchmarkLMDBGetNext(b *testing.B) {
	benchmarkGetNext(b, database.NewLMDB, nil, false)

	return
}

func BenchmarkLMDBGetNextInDB(b *testing.B) {
	benchmarkGetNext(b, database.NewLMDB, keyspace, false)

	return
}

func BenchmarkBoltPut(b *testing.B) {
	benchmarkPut(b, database.NewBoltDB, keyspace, false)

	return
}

func BenchmarkBoltGet(b *testing.B) {
	benchmarkGet(b, database.NewBoltDB, keyspace, false)

	return
}

func BenchmarkBoltGetNext(b *testing.B) {
	benchmarkGetNext(b, database.NewBoltDB, keyspace, false)

	return
}

func BenchmarkLevelPut(b *testing.B) {
	benchmarkPut(b, database.NewLevelDB, nil, false)

	return
}

func BenchmarkLevelGet(b *testing.B) {
	benchmarkGet(b, database.NewLevelDB, nil, false)

	return
}

func BenchmarkLevelGetNext(b *testing.B) {
	benchmarkGetNext(b, database.NewLevelDB, nil, false)

	return
}

func BenchmarkBadgerPut(b *testing.B) {
	benchmarkPut(b, database.NewBadgerDB, nil, false)

	return
}

func BenchmarkBadgerGet(b *testing.B) {
	benchmarkGet(b, database.NewBadgerDB, nil, false)

	return
}

func BenchmarkBadgerGetNext(b *testing.B) {
	benchmarkGetNext(b, database.NewBadgerDB, nil, false)

	return
}

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
