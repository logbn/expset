# Expiring Set

A generic set with time based eviction.

[![Go Reference](https://godoc.org/github.com/logbn/expset?status.svg)](https://godoc.org/github.com/logbn/expset)
[![License](https://img.shields.io/badge/License-Apache_2.0-dd6600.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/logbn/expset?4)](https://goreportcard.com/report/github.com/logbn/expset)
[![Go Coverage](https://github.com/logbn/expset/wiki/coverage.svg)](https://raw.githack.com/wiki/logbn/expset/coverage.html)

This set uses generics to contain any type of `comparable` element, evicting elements after a given expiration period.

## Usage

Create and start the set then add some items with time to live

```go
import "github.com/logbn/expset"

s := expset.New[string]()
s.Start()
defer s.Stop()

s.Add("test-1", 10 * time.Second)
s.Add("test-2", 20 * time.Second)
s.Add("test-3", 30 * time.Second)
```

Test whether set contains a value

```go
println(s.Contains("test-1"))
// output:
//   true

println(s.Contains("test-2"))
// output:
//   true
```

Observe expiration

```go
time.Sleep(10*time.Second)

println(s.Contains("test-1"))
// output:
//   false

println(s.Contains("test-2"))
// output:
//   true
```

Refresh an item to reset its expiration

```go
s.Refresh("test-2")

time.Sleep(10*time.Second)

println(s.Contains("test-2"))
// output:
//   true
```

## Concurrency

This package is thread safe.

## License

Expiring Set is licensed under Apache 2.0
