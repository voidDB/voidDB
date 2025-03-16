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
conducted on x86-64 and AArch64 instances.

### x86-64/`amd64`

#### AMD EPYC "Milan"

Amazon EC2 `r6a`, EBS `gp3`:

```txt
goos: linux
goarch: amd64
pkg: test
cpu: AMD EPYC 7R13 Processor
BenchmarkPopulateKeyVal-8   	  131072	     14407 ns/op
BenchmarkVoidPut-8          	  131072	     49785 ns/op
BenchmarkVoidGet-8          	  131072	      1371 ns/op
BenchmarkVoidGetNext-8      	  131072	       444.1 ns/op
BenchmarkLMDBPut-8          	  131072	     76957 ns/op
BenchmarkLMDBGet-8          	  131072	      1603 ns/op
BenchmarkLMDBGetNext-8      	  131072	       879.3 ns/op
BenchmarkBoltPut-8          	  131072	    229226 ns/op
BenchmarkBoltGet-8          	  131072	      2414 ns/op
BenchmarkBoltGetNext-8      	  131072	       452.2 ns/op
BenchmarkLevelPut-8         	  131072	    152578 ns/op
BenchmarkLevelGet-8         	  131072	     30373 ns/op
BenchmarkLevelGetNext-8     	  131072	      4196 ns/op
BenchmarkBadgerPut-8        	  131072	     89332 ns/op
BenchmarkBadgerGet-8        	  131072	     20421 ns/op
BenchmarkBadgerGetNext-8    	  131072	      2731 ns/op
BenchmarkNothing-8          	  131072	         0.2787 ns/op
```

> [R6a instances](https://aws.amazon.com/ec2/instance-types/r6a/) are powered
> by 3rd generation AMD EPYC processors ... and are an ideal fit for
> memory-intensive workloads, such as SQL and NoSQL databases; distributed web
> scale in-memory caches, such as Memcached and Redis; in-memory databases and
> real-time big data analytics, such as Apache Hadoop and Apache Spark
> clusters; and other enterprise applications.

#### Intel Xeon Platinum "Sapphire Rapids"

Amazon EC2 `r7i`, EBS `gp3`:

```txt
goos: linux
goarch: amd64
pkg: test
cpu: Intel(R) Xeon(R) Platinum 8488C
BenchmarkPopulateKeyVal-8   	  131072	     10750 ns/op
BenchmarkVoidPut-8          	  131072	     56492 ns/op
BenchmarkVoidGet-8          	  131072	      1227 ns/op
BenchmarkVoidGetNext-8      	  131072	       377.2 ns/op
BenchmarkLMDBPut-8          	  131072	     73205 ns/op
BenchmarkLMDBGet-8          	  131072	      1563 ns/op
BenchmarkLMDBGetNext-8      	  131072	       691.0 ns/op
BenchmarkBoltPut-8          	  131072	    417657 ns/op
BenchmarkBoltGet-8          	  131072	      2091 ns/op
BenchmarkBoltGetNext-8      	  131072	       271.7 ns/op
BenchmarkLevelPut-8         	  131072	     97302 ns/op
BenchmarkLevelGet-8         	  131072	     28205 ns/op
BenchmarkLevelGetNext-8     	  131072	      3799 ns/op
BenchmarkBadgerPut-8        	  131072	     90613 ns/op
BenchmarkBadgerGet-8        	  131072	     15375 ns/op
BenchmarkBadgerGetNext-8    	  131072	      3033 ns/op
BenchmarkNothing-8          	  131072	         0.3990 ns/op
```

> [R7i instances](https://aws.amazon.com/ec2/instance-types/r7i/) are ...
> powered by custom 4th Generation Intel Xeon Scalable processors (code named
> Sapphire Rapids) ... and ideal for all memory-intensive workloads (SQL and
> NoSQL databases), distributed web scale in-memory caches (Memcached and
> Redis), in-memory databases (SAP HANA), and real-time big data analytics
> (Apache Hadoop and Apache Spark clusters).

### AArch64/`arm64`

#### Ampere Altra

Google Cloud Compute Engine `t2a`, SSD persistent disk:

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-8   	  131072	     13850 ns/op
BenchmarkVoidPut-8          	  131072	     41688 ns/op
BenchmarkVoidGet-8          	  131072	      1803 ns/op
BenchmarkVoidGetNext-8      	  131072	       460.4 ns/op
BenchmarkLMDBPut-8          	  131072	     48451 ns/op
BenchmarkLMDBGet-8          	  131072	      2267 ns/op
BenchmarkLMDBGetNext-8      	  131072	       892.3 ns/op
BenchmarkBoltPut-8          	  131072	    192896 ns/op
BenchmarkBoltGet-8          	  131072	      3566 ns/op
BenchmarkBoltGetNext-8      	  131072	       451.0 ns/op
BenchmarkLevelPut-8         	  131072	     59104 ns/op
BenchmarkLevelGet-8         	  131072	     42006 ns/op
BenchmarkLevelGetNext-8     	  131072	      4860 ns/op
BenchmarkBadgerPut-8        	  131072	     47798 ns/op
BenchmarkBadgerGet-8        	  131072	     27400 ns/op
BenchmarkBadgerGetNext-8    	  131072	      3127 ns/op
BenchmarkNothing-8          	  131072	         0.4184 ns/op
```

#### Apple M1 Pro chip

Ubuntu VM in Multipass for macOS, MacBook Pro, Apple NVMe SSD:

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-2   	  131072	      9069 ns/op
BenchmarkVoidPut-2          	  131072	     14520 ns/op
BenchmarkVoidGet-2          	  131072	      1018 ns/op
BenchmarkVoidGetNext-2      	  131072	       243.0 ns/op
BenchmarkLMDBPut-2          	  131072	     24132 ns/op
BenchmarkLMDBGet-2          	  131072	      1455 ns/op
BenchmarkLMDBGetNext-2      	  131072	       583.5 ns/op
BenchmarkBoltPut-2          	  131072	     75558 ns/op
BenchmarkBoltGet-2          	  131072	      2089 ns/op
BenchmarkBoltGetNext-2      	  131072	       243.0 ns/op
BenchmarkLevelPut-2         	  131072	     40225 ns/op
BenchmarkLevelGet-2         	  131072	     27412 ns/op
BenchmarkLevelGetNext-2     	  131072	      2870 ns/op
BenchmarkBadgerPut-2        	  131072	     15450 ns/op
BenchmarkBadgerGet-2        	  131072	     21309 ns/op
BenchmarkBadgerGetNext-2    	  131072	     13626 ns/op
BenchmarkNothing-2          	  131072	         0.3249 ns/op
```

Native macOS, MacBook Pro, Apple NVMe SSD:

```txt
goos: darwin
goarch: arm64
pkg: test
cpu: Apple M1 Pro
BenchmarkPopulateKeyVal-10    	  131072	      4598 ns/op
BenchmarkVoidPut-10           	  131072	     10755 ns/op
BenchmarkVoidGet-10           	  131072	       852.8 ns/op
BenchmarkVoidGetNext-10       	  131072	       224.9 ns/op
BenchmarkLMDBPut-10           	  131072	      6280 ns/op
BenchmarkLMDBGet-10           	  131072	      1707 ns/op
BenchmarkLMDBGetNext-10       	  131072	       754.7 ns/op
BenchmarkBoltPut-10           	  131072	     61757 ns/op
BenchmarkBoltGet-10           	  131072	      1807 ns/op
BenchmarkBoltGetNext-10       	  131072	       439.5 ns/op
BenchmarkLevelPut-10          	  131072	     97582 ns/op
BenchmarkLevelGet-10          	  131072	     40275 ns/op
BenchmarkLevelGetNext-10      	  131072	      3687 ns/op
BenchmarkBadgerPut-10         	  131072	      7126 ns/op
BenchmarkBadgerGet-10         	  131072	     13894 ns/op
BenchmarkBadgerGetNext-10     	  131072	      2622 ns/op
BenchmarkNothing-10           	  131072	         0.3366 ns/op
```

#### AWS Graviton4

Amazon EC2 `r8g`, EBS `gp3`:

```txt
goos: linux
goarch: arm64
pkg: test
BenchmarkPopulateKeyVal-8   	  131072	     10807 ns/op
BenchmarkVoidPut-8          	  131072	     49006 ns/op
BenchmarkVoidGet-8          	  131072	      1225 ns/op
BenchmarkVoidGetNext-8      	  131072	       422.0 ns/op
BenchmarkLMDBPut-8          	  131072	     71656 ns/op
BenchmarkLMDBGet-8          	  131072	      1470 ns/op
BenchmarkLMDBGetNext-8      	  131072	       749.5 ns/op
BenchmarkBoltPut-8          	  131072	    119710 ns/op
BenchmarkBoltGet-8          	  131072	      2158 ns/op
BenchmarkBoltGetNext-8      	  131072	       354.1 ns/op
BenchmarkLevelPut-8         	  131072	     77564 ns/op
BenchmarkLevelGet-8         	  131072	     27252 ns/op
BenchmarkLevelGetNext-8     	  131072	      3379 ns/op
BenchmarkBadgerPut-8        	  131072	     90776 ns/op
BenchmarkBadgerGet-8        	  131072	     15028 ns/op
BenchmarkBadgerGetNext-8    	  131072	      1746 ns/op
BenchmarkNothing-8          	  131072	         0.3586 ns/op
```

> [R8g instances](https://aws.amazon.com/ec2/instance-types/r8g/), powered by
> the latest-generation AWS Graviton4 processors, ... are ideal for
> memory-intensive workloads, such as databases, in-memory caches, and
> real-time big data analytics.

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
