# voidDB

<a href="https://pkg.go.dev/github.com/voidDB/voidDB">
  <img src="https://pkg.go.dev/badge/github.com/voidDB/voidDB.svg" />
</a>
<div align="center">
  <img src="https://github.com/voidDB.png" width="230" />
</div>

voidDB is a [memory-mapped](https://man7.org/linux/man-pages/man2/mmap.2.html)
key-value store: simultaneously in-memory and persistent on disk. An embedded
database manager, it is meant to be integrated into application software to
eliminate protocol overheads and achieve zero-copy performance. This library
supplies interfaces for storage and retrieval of arbitrary bytes on 64-bit
computers running Linux and macOS.

voidDB features Put, Get, and Del operations as well as forward and backward
iteration over self-sorting data in ACID (atomic, consistent, isolated, and
durable) transactions. Readers retain a consistent view of the data throughout
their lifetime, even as newer transactions are being committed: only pages
freed by transactions older than the oldest surviving reader are actively
recycled.

voidDB employs a copy-on-write strategy to maintain data in a multi-version
concurrency-controlled (MVCC) B+ tree structure. It allows virtually any number
of concurrent readers, but only one active writer at any given moment. Readers
(and the sole writer) neither compete nor block one another, even though they
may originate from and operate within different threads and processes.

voidDB is resilient against torn writes. It automatically restores a database
to its last stable state in the event of a mid-write crash. Once a transaction
is committed and flushed to disk it is safe, but even if not it could do no
harm to existing data in storage. Applications need not be concerned about
broken lockfiles or lingering effects of unfinished transactions should an
uncontrolled shutdown occur; its design guarantees automatic and immediate
release of resources upon process termination.

## Benchmarks

voidDB consistently outperforms well-known key-value stores available to Go
developers based on B+ trees (LMDB, bbolt) and log-structured merge(LSM)-trees
(LevelDB, BadgerDB).

### 2,048 values × 256 KiB = 512 MiB

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-2          	    2048	    486747 ns/op
BenchmarkVoidPut-2                 	    2048	    242955 ns/op
BenchmarkVoidPutInKeyspace-2       	    2048	    252459 ns/op
BenchmarkVoidGet-2                 	    2048	       965.0 ns/op
BenchmarkVoidGetInKeyspace-2       	    2048	       871.5 ns/op
BenchmarkVoidGetNext-2             	    2048	       704.4 ns/op
BenchmarkVoidGetNextInKeyspace-2   	    2048	       611.5 ns/op
BenchmarkLMDBPut-2                 	    2048	    501073 ns/op
BenchmarkLMDBPutInDB-2             	    2048	    260418 ns/op
BenchmarkLMDBGet-2                 	    2048	      2163 ns/op
BenchmarkLMDBGetInDB-2             	    2048	      2159 ns/op
BenchmarkLMDBGetNext-2             	    2048	      1920 ns/op
BenchmarkLMDBGetNextInDB-2         	    2048	      2004 ns/op
BenchmarkBoltPut-2                 	    2048	    483040 ns/op
BenchmarkBoltGet-2                 	    2048	      2412 ns/op
BenchmarkBoltGetNext-2             	    2048	       856.4 ns/op
BenchmarkLevelPut-2                	    2048	    902054 ns/op
BenchmarkLevelGet-2                	    2048	     92537 ns/op
BenchmarkLevelGetNext-2            	    2048	     85487 ns/op
BenchmarkBadgerPut-2               	    2048	    632082 ns/op
BenchmarkBadgerGet-2               	    2048	     32034 ns/op
BenchmarkBadgerGetNext-2           	    2048	     31033 ns/op
BenchmarkNothing-2                 	    2048	         0.4678 ns/op
```

### 32,768 values × 16 KiB = 512 MiB

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-2          	   32768	     32403 ns/op
BenchmarkVoidPut-2                 	   32768	     23288 ns/op
BenchmarkVoidPutInKeyspace-2       	   32768	     22056 ns/op
BenchmarkVoidGet-2                 	   32768	       986.5 ns/op
BenchmarkVoidGetInKeyspace-2       	   32768	       981.5 ns/op
BenchmarkVoidGetNext-2             	   32768	       551.0 ns/op
BenchmarkVoidGetNextInKeyspace-2   	   32768	       381.2 ns/op
BenchmarkLMDBPut-2                 	   32768	     36842 ns/op
BenchmarkLMDBPutInDB-2             	   32768	     24648 ns/op
BenchmarkLMDBGet-2                 	   32768	      1401 ns/op
BenchmarkLMDBGetInDB-2             	   32768	      1546 ns/op
BenchmarkLMDBGetNext-2             	   32768	       827.4 ns/op
BenchmarkLMDBGetNextInDB-2         	   32768	       838.8 ns/op
BenchmarkBoltPut-2                 	   32768	     72147 ns/op
BenchmarkBoltGet-2                 	   32768	      2177 ns/op
BenchmarkBoltGetNext-2             	   32768	       445.6 ns/op
BenchmarkLevelPut-2                	   32768	     69430 ns/op
BenchmarkLevelGet-2                	   32768	     21170 ns/op
BenchmarkLevelGetNext-2            	   32768	     11026 ns/op
BenchmarkBadgerPut-2               	   32768	     42966 ns/op
BenchmarkBadgerGet-2               	   32768	     25837 ns/op
BenchmarkBadgerGetNext-2           	   32768	      4509 ns/op
BenchmarkNothing-2                 	   32768	         0.3255 ns/op
```

### 524,288 values × 1 KiB = 512 MiB

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-2          	  524288	      3623 ns/op
BenchmarkVoidPut-2                 	  524288	      8090 ns/op
BenchmarkVoidPutInKeyspace-2       	  524288	      8122 ns/op
BenchmarkVoidGet-2                 	  524288	      1321 ns/op
BenchmarkVoidGetInKeyspace-2       	  524288	      1338 ns/op
BenchmarkVoidGetNext-2             	  524288	       179.0 ns/op
BenchmarkVoidGetNextInKeyspace-2   	  524288	       166.9 ns/op
BenchmarkLMDBPut-2                 	  524288	     15095 ns/op
BenchmarkLMDBPutInDB-2             	  524288	     11524 ns/op
BenchmarkLMDBGet-2                 	  524288	      1980 ns/op
BenchmarkLMDBGetInDB-2             	  524288	      1977 ns/op
BenchmarkLMDBGetNext-2             	  524288	       595.8 ns/op
BenchmarkLMDBGetNextInDB-2         	  524288	       623.5 ns/op
BenchmarkBoltPut-2                 	  524288	    205160 ns/op
BenchmarkBoltGet-2                 	  524288	      2381 ns/op
BenchmarkBoltGetNext-2             	  524288	       184.1 ns/op
BenchmarkLevelPut-2                	  524288	     18287 ns/op
BenchmarkLevelGet-2                	  524288	     20767 ns/op
BenchmarkLevelGetNext-2            	  524288	      1107 ns/op
BenchmarkBadgerPut-2               	  524288	      6386 ns/op
BenchmarkBadgerGet-2               	  524288	     38710 ns/op
BenchmarkBadgerGetNext-2           	  524288	     16545 ns/op
BenchmarkNothing-2                 	  524288	         0.3314 ns/op
```

## Getting Started

[Install Go](https://go.dev/doc/install) to begin developing with voidDB.

```bash
$ go version
go version go1.24.0 linux/arm64
```

Then, import voidDB in your Go application. The following would result in the
creation of a database file and its reader table in the working directory. Set
the database capacity to any reasonably large value to make sufficient room for
the data you intend to store, even if it exceeds the total amount of physical
memory; neither memory nor disk is immediately consumed to capacity.

```go
package main

import (
	"errors"
	"os"

	"github.com/voidDB/voidDB"
)

func main() {
	const (
		capacity = 1 << 40 // 1 TiB
		path     = "void"
	)

	void, err := voidDB.NewVoid(path, capacity)

	if errors.Is(err, os.ErrExist) {
		void, err = voidDB.OpenVoid(path, capacity)
	}

	if err != nil {
		panic(err)
	}

	defer void.Close()
}
```

Use `*Void.View` (or `*Void.Update` only when modifying data) for convenience
and peace of mind. Ensure all changes are safely synced to disk with `mustSync`
set to `true` if even the slightest risk of losing those changes is a concern.

```go
mustSync := true

err = void.Update(mustSync,
	func(txn *voidDB.Txn) error {
		return txn.Put(
			[]byte("greeting"),
			[]byte("Hello, World!"),
		)
	},
)
if err != nil {
	panic(err)
}
```

Open a cursor if more than one keyspace is required. An application can map
different values to the same key so long as they reside in separate keyspaces.
The transaction handle doubles as a cursor in the default keyspace.

```go
cur0, _ := txn.OpenCursor([]byte("hello"))

cur0.Put([]byte("greeting"),
	[]byte("Hello, World!"),
)

cur1, _ := txn.OpenCursor([]byte("goodbye"))

cur1.Put([]byte("greeting"),
	[]byte("さらばこの世、わらわはもう寝るぞよ。"),
)

if val, err := cur0.Get([]byte("greeting")); err == nil {
	log.Printf("%s", val) // Hello, World!
}

if val, err := cur1.Get([]byte("greeting")); err == nil {
	log.Printf("%s", val) // さらばこの世、わらわはもう寝るぞよ。
}
```

To iterate over a keyspace, use `*cursor.Cursor.GetNext`/`GetPrev`. Position
the cursor with `*cursor.Cursor.Get`/`GetFirst`/`GetLast`.

```go
for {
	key, val, err := cur.GetNext()

	if errors.Is(err, common.ErrorNotFound) {
		break
	}

	log.Printf("%s -> %s", key, val)
}
```

## Author

voidDB builds upon ideas in the celebrated [Lightning Memory-Mapped Database
Manager](http://www.lmdb.tech/doc/) on several key points of its high-level
design, but otherwise it is implemented from scratch to break free of
limitations in function, performance, and clarity.

voidDB is a cherished toy, a journey into the Unknown, a heroic struggle, and a
work of love. It is the “Twee!” of a bird; a tree falling in the forest; yet
another programmer pouring their drop into the proverbial [bit] bucket. Above
all, it is a shrine unto simple, readable, and functional code; an assertion
that the dichotomy between such aesthetics and practical performance is mere
illusion.

Copyright 2024 Joel Ling
