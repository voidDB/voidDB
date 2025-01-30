# voidDB

<a href="https://pkg.go.dev/github.com/voidDB/voidDB">
  <img src="https://pkg.go.dev/badge/github.com/voidDB/voidDB.svg" />
</a>
<div align="center">
  <img src="https://github.com/voidDB.png" width="230" />
</div>

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
