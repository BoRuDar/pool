# pool
Memory pool based on `sync.Pool` to use with generic structs.
Main pros:
* easy to use - just 3 lines of code
* automatic reset of the used objects before returning back to the pool
* no need to manually cast and check types
* safe for a concurrent access

## Quick start

1. Create pool
2. Get an item from the pool
3. Return the item to the pool
```go
    pool := New[userTestStr]()
    u := pool.Get()
    defer pool.Return(u)
```

## Benchmarks 

```
go test -bench=. -benchmem -count 3 -benchtime 5s
goos: linux
goarch: amd64
pkg: github.com/BoRuDar/pool
cpu: AMD Ryzen 9 7950X3D 16-Core Processor          
BenchmarkNoPool-32               5500141               952.9 ns/op           648 B/op          8 allocs/op
BenchmarkNoPool-32               6639531               977.9 ns/op           648 B/op          8 allocs/op
BenchmarkNoPool-32               5964730              1013 ns/op             648 B/op          8 allocs/op
BenchmarkWithPool-32             6285232               924.7 ns/op           360 B/op          4 allocs/op
BenchmarkWithPool-32             6162657               990.5 ns/op           361 B/op          4 allocs/op
BenchmarkWithPool-32             6403029               974.1 ns/op           361 B/op          4 allocs/op
PASS
ok      github.com/BoRuDar/pool 41.890s
```