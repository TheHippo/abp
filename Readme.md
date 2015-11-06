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

## License

See `License` file
