//go:generate protoc -I=internal/proto -I=vendor --gogofaster_out=internal/model internal/proto/update.proto

package copydb
