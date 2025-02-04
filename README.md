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
computers running the Linux operating system.

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

```txt
goos: linux
goarch: arm64
pkg: github.com/voidDB/voidDB/test
BenchmarkVoidPut-2         	  131072	     14933 ns/op
BenchmarkVoidGet-2         	  131072	      1060 ns/op
BenchmarkVoidGetNext-2     	  131072	       245.8 ns/op
BenchmarkLMDBPut-2         	  131072	     22414 ns/op
BenchmarkLMDBGet-2         	  131072	      1826 ns/op
BenchmarkLMDBGetNext-2     	  131072	       602.2 ns/op
BenchmarkBoltPut-2         	  131072	     66984 ns/op
BenchmarkBoltGet-2         	  131072	      2552 ns/op
BenchmarkBoltGetNext-2     	  131072	       254.6 ns/op
BenchmarkLevelPut-2        	  131072	     44182 ns/op
BenchmarkLevelGet-2        	  131072	     30949 ns/op
BenchmarkLevelGetNext-2    	  131072	      3441 ns/op
BenchmarkBadgerPut-2       	  131072	     15182 ns/op
BenchmarkBadgerGet-2       	  131072	     33114 ns/op
BenchmarkBadgerGetNext-2   	  131072	     12895 ns/op
BenchmarkNothing-2         	  131072	         0.3239 ns/op
```

## Getting Started

To begin developing with voidDB, [Install Go](https://go.dev/doc/install) on
your 64-bit Linux machine.

```bash
$ go version
go version go1.22.3 linux/arm64
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
and peace of mind. Ensure all changes are safely synced to disk with mustSync
set to true if even the slightest risk of losing those changes is a concern.

```go
mustSync := true

err = void.Update(mustSync,
	func(txn *voidDB.Txn) (err error) {
		err = txn.Put(
			[]byte("greeting"),
			[]byte("Hello, World!"),
		)

		return
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
cur0, err := txn.OpenCursor(
	[]byte("hello"),
)

err = cur0.Put(
	[]byte("greeting"),
	[]byte("Hello, World!"),
)

cur1, err := txn.OpenCursor(
	[]byte("goodbye"),
)

err = cur1.Put(
	[]byte("greeting"),
	[]byte("さようなら、世界。"),
)

val, err := cur0.Get(
	[]byte("greeting"),
)

log.Printf("%s", val) // Hello, World!

val, err = cur1.Get(
	[]byte("greeting"),
)

log.Printf("%s", val) // さようなら、世界。
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
