//go:generate protoc -I=proto -I=vendor --gogofaster_out=.  proto/update.proto

//go:generate easyjson -all examples/providers/provider.go examples/providers/request.go

package copydb
