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
BenchmarkVoidPut-2         	  131072	     15326 ns/op
BenchmarkVoidGet-2         	  131072	      1046 ns/op
BenchmarkVoidGetNext-2     	  131072	       235.6 ns/op
BenchmarkLMDBPut-2         	  131072	     25671 ns/op
BenchmarkLMDBGet-2         	  131072	      1478 ns/op
BenchmarkLMDBGetNext-2     	  131072	       601.5 ns/op
BenchmarkBoltPut-2         	  131072	     70759 ns/op
BenchmarkBoltGet-2         	  131072	      2789 ns/op
BenchmarkBoltGetNext-2     	  131072	       246.9 ns/op
BenchmarkLevelPut-2        	  131072	     46364 ns/op
BenchmarkLevelGet-2        	  131072	     29708 ns/op
BenchmarkLevelGetNext-2    	  131072	      3331 ns/op
BenchmarkBadgerPut-2       	  131072	     15314 ns/op
BenchmarkBadgerGet-2       	  131072	     21152 ns/op
BenchmarkBadgerGetNext-2   	  131072	      2019 ns/op
BenchmarkNothing-2         	  131072	         0.3204 ns/op
```
