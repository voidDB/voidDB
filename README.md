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
conducted on x86-64 and AArch64 instances hosted on Google Cloud (N2, T2A/D
machine series, 8 vCPUs, 32 GB memory).

### Put

#### 4,096 × 256-KiB random values

```txt
|                   (ms/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      1.16 |         1.12 |               1.16 |
| voidDB (named keyspace)   |      1.16 |         1.12 |               1.16 |
| LMDB (default keyspace)   |      1.90 |         2.00 |               1.93 |
| Bolt                      |      1.96 |         1.76 |               2.45 |
| LMDB (named keyspace)     |      2.11 |         2.26 |               2.15 |
| LevelDB                   |      2.37 |         2.23 |               2.75 |
| BadgerDB                  |      3.22 |         5.31 |               3.14 |
```

#### 65,536 × 16-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (named keyspace)   |      82.7 |         82.8 |               85.0 |
| voidDB (default keyspace) |      83.5 |         77.7 |               85.3 |
| LMDB (named keyspace)     |     152   |        154   |              151   |
| LMDB (default keyspace)   |     157   |        156   |              157   |
| BadgerDB                  |     195   |        225   |              183   |
| Bolt                      |     244   |        217   |              429   |
| LevelDB                   |     362   |        310   |              303   |
```

#### 1,048,576 × 1-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      18.7 |         19.5 |               21.4 |
| voidDB (named keyspace)   |      18.4 |         19.8 |               21.9 |
| BadgerDB                  |      24.2 |         23.9 |               25.8 |
| LMDB (named keyspace)     |      28.2 |         32.3 |               33.5 |
| LMDB (default keyspace)   |      29.5 |         34.4 |               36.6 |
| LevelDB                   |      72.5 |         58.8 |              165   |
| Bolt                      | timed out | timed out    | timed out          |
```

### Get

#### 4,096 × 256-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (named keyspace)   |      1.64 |         1.76 |               1.66 |
| voidDB (default keyspace) |      1.65 |         1.74 |               1.78 |
| LMDB (named keyspace)     |      4.97 |         3.97 |               4.46 |
| LMDB (default keyspace)   |      5.00 |         4.03 |               4.33 |
| Bolt                      |      5.88 |         5.17 |               5.76 |
| LevelDB                   |    115    |       125    |             224    |
| BadgerDB                  |    301    |       142    |             618    |
```

#### 65,536 × 16-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      1.73 |         2.01 |               2.29 |
| voidDB (named keyspace)   |      1.76 |         2.02 |               2.11 |
| LMDB (named keyspace)     |      2.47 |         2.73 |               3.05 |
| LMDB (default keyspace)   |      2.62 |         2.60 |               3.12 |
| Bolt                      |      3.68 |         3.87 |               4.59 |
| LevelDB                   |     27.4  |        34.6  |              46.1  |
| BadgerDB                  |     21.1  |        41.8  |              71.2  |
```

#### 1,048,576 × 1-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      1.76 |         2.15 |               2.54 |
| voidDB (named keyspace)   |      2.07 |         2.58 |               3.00 |
| LMDB (named keyspace)     |      2.09 |         2.68 |               2.97 |
| LMDB (default keyspace)   |      2.22 |         2.68 |               3.02 |
| BadgerDB                  |     23.8  |        22.8  |              31.1  |
| LevelDB                   |     27.3  |        45.7  |              40.7  |
| Bolt                      | timed out | timed out    | timed out          |
```

### GetNext

#### 4,096 × 256-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      1.17 |         .938 |               1.15 |
| voidDB (named keyspace)   |      1.27 |         .939 |               1.12 |
| Bolt                      |      2.13 |        1.55  |               1.83 |
| LMDB (named keyspace)     |      4.54 |        3.50  |               3.80 |
| LMDB (default keyspace)   |      4.65 |        3.44  |               3.90 |
| LevelDB                   |    107    |      110     |             181    |
| BadgerDB                  |    198    |       58.7   |             415    |
```

#### 65,536 × 16-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      .808 |         .721 |               .891 |
| voidDB (named keyspace)   |      .849 |         .730 |               .874 |
| Bolt                      |     1.18  |         .869 |               .919 |
| LMDB (default keyspace)   |     1.78  |        1.47  |              1.64  |
| LMDB (named keyspace)     |     1.78  |        1.44  |              1.68  |
| LevelDB                   |    14.7   |       14.7   |             22.3   |
| BadgerDB                  |    26.1   |        6.54  |             24.1   |
```

#### 1,048,576 × 1-KiB random values

```txt
|                   (μs/op) | AMD Milan | Ampere Altra | Intel Cascade Lake |
| ------------------------- | --------- | ------------ | ------------------ |
| voidDB (default keyspace) |      .307 |         .270 |               .384 |
| voidDB (named keyspace)   |      .296 |         .515 |               .402 |
| LMDB (named keyspace)     |      .615 |         .515 |               .637 |
| LMDB (default keyspace)   |      .607 |         .519 |               .642 |
| LevelDB                   |     1.38  |        1.80  |              1.92  |
| BadgerDB                  |     4.03  |        5.60  |             15.0   |
| Bolt                      | timed out | timed out    | timed out          |
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
