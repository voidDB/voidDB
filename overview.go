// voidDB is a [memory-mapped] key-value store: simultaneously in-memory and
// persistent on disk. An embedded database manager, it is meant to be
// integrated into application software to eliminate protocol overheads and
// achieve zero-copy performance. This library supplies interfaces for storage
// and retrieval of arbitrary bytes on 64-bit computers running the Linux
// operating system.
//
// voidDB features Put, Get, and Del operations as well as forward and backward
// iteration over self-sorting data in ACID (atomic, consistent, isolated, and
// durable) transactions. Readers retain a consistent view of the data
// throughout their lifetime, even as newer transactions are being committed:
// only pages freed by transactions older than the oldest surviving reader are
// actively recycled.
//
// voidDB employs a copy-on-write strategy to maintain data in a multi-version
// concurrency-controlled (MVCC) B+ tree structure. It allows virtually any
// number of concurrent readers, but only one active writer at any given
// moment. Readers (and the sole writer) neither compete nor block one another,
// even though they may originate from and operate within different threads and
// processes.
//
// voidDB applications need not be concerned about broken lockfiles or
// lingering effects of unfinished transactions should an uncontrolled shutdown
// occur; its design guarantees automatic and immediate release of resources
// upon process termination.
//
// # Getting Started
//
// TODO.
//
// # Author
//
// voidDB is a cherished toy, a journey into the Unknown, a heroic struggle,
// and a work of love. It is the “Twee!” of a bird; a tree falling in the
// forest; yet another programmer pouring their drop into the proverbial [bit]
// bucket. Above all, it is a shrine unto simple, readable, and functional
// code; an assertion that the dichotomy between such aesthetics and practical
// performance is mere illusion.
//
// Copyright 2024 Joel Ling.
//
// [memory-mapped]: https://man7.org/linux/man-pages/man2/mmap.2.html
package voidDB
