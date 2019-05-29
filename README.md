# CopyDB

[![Build Status](https://travis-ci.org/regeda/copydb.svg?branch=master)](https://travis-ci.org/regeda/copydb)
[![codecov](https://codecov.io/gh/regeda/copydb/branch/master/graph/badge.svg)](https://codecov.io/gh/regeda/copydb)
[![Go Report Card](https://goreportcard.com/badge/github.com/regeda/copydb)](https://goreportcard.com/report/github.com/regeda/copydb)
[![GoDoc](https://godoc.org/github.com/regeda/copydb?status.svg)](https://godoc.org/github.com/regeda/copydb)

CopyDB is **in-memory** database replicated through Redis pub/sub approach. All writes are synced with Redis and then replicated to all CopyDB instances.
CopyDB is an **embedded** database and it gives you a freedom to organize any cache patterns inside your application.
CopyDB is the best choice for data that updates frequent to provide **eventual consistency**.

## How it works

CopyDB loads all available items in memory from Redis after start. Then the database subscribes on items updates.

> Each item contains a version to guarantee a consistent update on a particular instance.
> An item will be reloaded from Redis if a new version is not greater than existed by 1.

The database supports the following operations: SET, UNSET and REMOVE. It's described by `Item` interface:
```go
type Item interface {
	Set(name string, data []byte)
	Unset(name string)
	Remove()
}
```

Read queries should be performed using embedded `Query*` functions:
- QueryByID
- QueryAll
- QueryPool
- QueryStats

> Queries execution blocks items updates. Take the note you must write optimized queries.

## Monitoring

CopyDB exposes essential metrics to get health status.
```go
type Stats struct {
	ItemsApplied           int
	ItemsFailed            int
	ItemsEvicted           int
	VersionConfictDetected int
	DBScanned              int
}
```
