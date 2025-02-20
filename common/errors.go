package common

import (
	"errors"
)

var (
	ErrorCorrupt  = errors.New("voidDB: encountered a page of unexpected type")
	ErrorDeleted  = errors.New("voidDB: record deleted no longer accessible")
	ErrorFull     = errors.New("voidDB: write would exceed scope of memory map")
	ErrorInvalid  = errors.New("voidDB: database file format not recognised")
	ErrorNotFound = errors.New("voidDB: record not found by key or cursor")
	ErrorResized  = errors.New("voidDB: database file larger than memory map")
)
