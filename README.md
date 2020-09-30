[![Go Report Card](https://goreportcard.com/badge/github.com/liangyaopei/goroutineinfo)](https://goreportcard.com/report/github.com/liangyaopei/goroutineinfo)
[![GoDoc](https://godoc.org/github.com/liangyaopei/goroutineinfo?status.svg)](http://godoc.org/github.com/liangyaopei/goroutineinfo)

A repository provides access to goroutine's ID, state and otehr information by parsing message invoking ` runtime.Stack`

## Install
```
go get -u github.com/liangyaopei/goroutineinfo
```

## Example
```go
func TestGetInfoSingle(t *testing.T) {
	stacks := goroutineinfo.GetInfo(false)
	for _, stack := range stacks {
		t.Logf("id:[%d],state:[%s]", stack.ID(), stack.State())
	}
}
```

## Benchmark
```go
func BenchmarkGetInfoSingle(b *testing.B) {
	_ = goroutineinfo.GetInfo(false)
}
```
testing
```
go test ./... -bench=BenchmarkGetInfoSingle -benchmem -run=^$ -count=10
```
time cost is  0.000034 ns/op average.

