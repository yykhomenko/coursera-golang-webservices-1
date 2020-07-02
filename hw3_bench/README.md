Setup:
```bash
go get -u github.com/mailru/easyjson/...
go install github.com/mailru/easyjson
```

```bash
easyjson -all fast.go
```
error if package named "main" - rename and generate

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
BenchmarkSlow-8               48          25039272 ns/op        18968737 B/op     195837 allocs/op
BenchmarkFast-8              691           1713600 ns/op          560695 B/op       7513 allocs/op
PASS
ok      github.com/yykhomenko/coursera-golang-webservices-1/hw3_bench   2.620s
```
