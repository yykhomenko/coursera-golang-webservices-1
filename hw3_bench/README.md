#hw3_bench

Program optimization.

https://www.coursera.org/learn/golang-webservices-1/home/welcome

The task completed.

Setup:
```bash
go get -u github.com/mailru/easyjson/...
go install github.com/mailru/easyjson
cd hw3_bench
easyjson -all fast.go
```
easyjson throws error if package named "main" - rename package and generate again.

Commands:
```bash
go test -v
go test -bench=Fast -benchmem -cpuprofile=cpu.out -memprofile=mem.out
go tool pprof hw3_bench.test cpu.out
go tool pprof hw3_bench.test mem.out
```

Results:
```bash
mac:hw3_bench yykhomenko$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/yykhomenko/coursera-golang-webservices-1/hw3_bench
BenchmarkSlow-8               48          24930166 ns/op        18978971 B/op     195840 allocs/op
BenchmarkFast-8              934           1296906 ns/op          301110 B/op       3985 allocs/op
PASS
ok      github.com/yykhomenko/coursera-golang-webservices-1/hw3_bench   2.620s
```
