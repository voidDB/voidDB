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

voidDB outperforms well-known key-value stores available to Go developers that
are based on B+ trees (LMDB, bbolt) and log-structured merge(LSM)-trees
(LevelDB, BadgerDB), in [preliminary performance tests](test/bench_test.go)
conducted on x86-64 and AArch64 instances hosted on Google Cloud (T2A/D machine
series).

### Put

#### 4,096 × 256-KiB random values

```txt
|                           | AMD Milan (ms/op) | Ampere Altra (ms/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (named keyspace)   |              1.24 |                 1.20 |
| voidDB (default keyspace) |              1.25 |                 1.20 |
| Bolt                      |              2.07 |                 1.83 |
| LMDB (named keyspace)     |              2.18 |                 2.14 |
| LMDB (default keyspace)   |              2.34 |                 2.52 |
| BadgerDB                  |              3.30 |                 3.63 |
| LevelDB                   |              3.74 |                 3.22 |
```

#### 65,536 × 16-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (named keyspace)   |              89.3 |                 87.8 |
| voidDB (default keyspace) |              89.4 |                 88.2 |
| LMDB (named keyspace)     |             154   |                136   |
| LMDB (default keyspace)   |             195   |                194   |
| BadgerDB                  |             214   |                225   |
| Bolt                      |             244   |                218   |
| LevelDB                   |             273   |                227   |
```

#### 1,048,576 × 1-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (default keyspace) |              21.1 |                 22.1 |
| voidDB (named keyspace)   |              20.9 |                 22.9 |
| LMDB (default keyspace)   |              36.0 |                 43.4 |
| LMDB (named keyspace)     |              36.3 |                 42.1 |
| LevelDB                   |              66.8 |                 56.4 |
| Bolt                      | (timed out)       | (timed out)          |
| BadgerDB                  | (crashed)         | (crashed)            |
```

### Get

#### 4,096 × 256-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (named keyspace)   |              1.69 |                 1.72 |
| voidDB (default keyspace) |              1.84 |                 1.74 |
| LMDB (default keyspace)   |              4.95 |                 4.11 |
| LMDB (named keyspace)     |              5.13 |                 4.28 |
| Bolt                      |              5.71 |                 5.31 |
| LevelDB                   |            121    |               223    |
| BadgerDB                  |            261    |               128    |
```

#### 65,536 × 16-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (default keyspace) |              1.82 |                 1.98 |
| voidDB (named keyspace)   |              2.02 |                 2.06 |
| LMDB (named keyspace)     |              2.72 |                 2.67 |
| LMDB (default keyspace)   |              2.81 |                 2.74 |
| Bolt                      |              3.66 |                 4.17 |
| LevelDB                   |             28.3  |                34.7  |
| BadgerDB                  |             83.2  |                46.3  |
```

#### 1,048,576 × 1-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| LMDB (default keyspace)   |              2.33 |                 2.65 |
| LMDB (named keyspace)     |              2.44 |                 2.67 |
| voidDB (named keyspace)   |              2.68 |                 2.12 |
| voidDB (default keyspace) |              3.30 |                 2.77 |
| LevelDB                   |             30.2  |                45.0  |
| Bolt                      | (timed out)       | (timed out)          |
| BadgerDB                  | (crashed)         | (crashed)            |
```

### GetNext

#### 4,096 × 256-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (named keyspace)   |              1.20 |                 .995 |
| voidDB (default keyspace) |              1.18 |                1.02  |
| Bolt                      |              2.30 |                1.74  |
| LMDB (named keyspace)     |              4.71 |                3.71  |
| LMDB (default keyspace)   |              4.73 |                3.77  |
| LevelDB                   |            104    |              114     |
| BadgerDB                  |            257    |               72.1   |
```

#### 65,536 × 16-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| voidDB (default keyspace) |              .922 |                 .733 |
| voidDB (named keyspace)   |              .957 |                 .727 |
| Bolt                      |             1.20  |                 .954 |
| LMDB (default keyspace)   |             1.97  |                1.58  |
| LMDB (named keyspace)     |             1.98  |                1.55  |
| BadgerDB                  |            11.5   |               10.1   |
| LevelDB                   |            55.8   |               15.1   |
```

#### 1,048,576 × 1-KiB random values

```txt
|                           | AMD Milan (μs/op) | Ampere Altra (μs/op) |
| ------------------------- | ----------------- | -------------------- |
| LMDB (named keyspace)     |              .677 |                 .571 |
| LMDB (default keyspace)   |              .705 |                 .573 |
| voidDB (default keyspace) |              .845 |                 .280 |
| voidDB (named keyspace)   |             1.03  |                1.51  |
| LevelDB                   |            20.6   |              237     |
| Bolt                      | (timed out)       | (timed out)          |
| BadgerDB                  | (crashed)         | (crashed)            |
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
