# AveragedBufferPool [![GoDoc](https://godoc.org/github.com/TheHippo/abp?status.png)](https://godoc.org/github.com/TheHippo/abp) [![Build Status](https://travis-ci.org/TheHippo/abp.svg?branch=master)](https://travis-ci.org/TheHippo/abp) [![Build Status](https://drone.io/github.com/TheHippo/abp/status.png)](https://drone.io/github.com/TheHippo/abp/latest) [![Coverage Status](https://coveralls.io/repos/TheHippo/abp/badge.svg?branch=master&service=github)](https://coveralls.io/github/TheHippo/abp?branch=master)

`abp` implements a buffer pool based on fixed sized channels aka [leaky buffers](http://golang.org/doc/effective_go.html#leaky_buffer).

The average size of a buffer get calculated based on the sizes of the returned buffers and is calculated using a ring.

## Installation

`go get github.com/TheHippo/abp`

## Documentation

See [godoc.org](http://godoc.org/github.com/TheHippo/abp) or use `godoc github.com/TheHippo/abp`

## Example

```go
pool := abp.NewAveragedBufferPool(10, 10000, 100, 1.2)

buf := pool.Get()
// do something and than
pool.Put(buf)

// if do not neet any buffers anymore, close
pool.Close()
```

### Benchmark

```
BenchmarkSimplePool100-4         1000000          1245 ns/op           0 B/op          0 allocs/op
BenchmarkSimplePool1000-4        1000000          1299 ns/op           0 B/op          0 allocs/op
BenchmarkSimplePool10000-4       1000000          1296 ns/op           0 B/op          0 allocs/op
BenchmarkParallelPool100-4       1000000          1588 ns/op           0 B/op          0 allocs/op
BenchmarkParallelPool1000-4      1000000          1433 ns/op           0 B/op          0 allocs/op
BenchmarkParallelPool10000-4     1000000          1526 ns/op           0 B/op          0 allocs/op
BenchmarkMove1000-4             200000000            6.02 ns/op        0 B/op          0 allocs/op
BenchmarkMove10000-4            200000000            6.14 ns/op        0 B/op          0 allocs/op
BenchmarkMoveFull1000-4         300000000            5.56 ns/op        0 B/op          0 allocs/op
BenchmarkMoveFull10000-4        300000000            5.25 ns/op        0 B/op          0 allocs/op
BenchmarkMoveAndSet1000-4       200000000            6.07 ns/op        0 B/op          0 allocs/op
BenchmarkMoveAndSet10000-4      200000000            6.19 ns/op        0 B/op          0 allocs/op
```

## License

See `License` file
