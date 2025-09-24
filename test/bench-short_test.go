package test

import (
	"testing"

	// _ "github.com/ianlancetaylor/cgosymbolizer"

	"test/database"
)

func BenchmarkVoidPutShort(b *testing.B) {
	benchmarkPut(b, database.NewVoidDB, nil, true)

	return
}

func BenchmarkVoidPutInKeyspaceShort(b *testing.B) {
	benchmarkPut(b, database.NewVoidDB, keyspace, true)

	return
}

func BenchmarkVoidGetShort(b *testing.B) {
	benchmarkGet(b, database.NewVoidDB, nil, true)

	return
}

func BenchmarkVoidGetInKeyspaceShort(b *testing.B) {
	benchmarkGet(b, database.NewVoidDB, keyspace, true)

	return
}

func BenchmarkLMDBPutShort(b *testing.B) {
	benchmarkPut(b, database.NewLMDB, nil, true)

	return
}

func BenchmarkLMDBPutInDBShort(b *testing.B) {
	benchmarkPut(b, database.NewLMDB, keyspace, true)

	return
}

func BenchmarkLMDBGetShort(b *testing.B) {
	benchmarkGet(b, database.NewLMDB, nil, true)

	return
}

func BenchmarkLMDBGetInDBShort(b *testing.B) {
	benchmarkGet(b, database.NewLMDB, keyspace, true)

	return
}

func BenchmarkBoltPutShort(b *testing.B) {
	benchmarkPut(b, database.NewBoltDB, keyspace, true)

	return
}

func BenchmarkBoltGetShort(b *testing.B) {
	benchmarkGet(b, database.NewBoltDB, keyspace, true)

	return
}

func BenchmarkLevelPutShort(b *testing.B) {
	benchmarkPut(b, database.NewLevelDB, nil, true)

	return
}

func BenchmarkLevelGetShort(b *testing.B) {
	benchmarkGet(b, database.NewLevelDB, nil, true)

	return
}

func BenchmarkBadgerPutShort(b *testing.B) {
	benchmarkPut(b, database.NewBadgerDB, nil, true)

	return
}

func BenchmarkBadgerGetShort(b *testing.B) {
	benchmarkGet(b, database.NewBadgerDB, nil, true)

	return
}
