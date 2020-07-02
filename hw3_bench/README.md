Setup:
```bash
go get -u github.com/mailru/easyjson/...
go install github.com/mailru/easyjson
easyjson -all fast.go
```
error if package named "main" - rename and generate

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
BenchmarkSlow-8               46          25229523 ns/op        18956210 B/op     195836 allocs/op
BenchmarkFast-8              699           1710098 ns/op          560682 B/op       7513 allocs/op
PASS
ok      github.com/yykhomenko/coursera-golang-webservices-1/hw3_bench   2.620s
```
