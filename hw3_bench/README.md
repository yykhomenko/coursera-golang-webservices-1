Commands:
```bash
go get github.com/mailru/easyjson/...
go install github.com/mailru/easyjson
easyjson -all fast.go
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
BenchmarkSlow-8               46          25664575 ns/op        18978185 B/op     195841 allocs/op
BenchmarkFast-8              667           1734744 ns/op          560676 B/op       7513 allocs/op
PASS
ok      github.com/yykhomenko/coursera-golang-webservices-1/hw3_bench   2.620s
```
