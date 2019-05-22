.PHONY: all clean generate build

all: clean generate build

clean:
	rm -rf bin

build:
	GOGC=off go build -o bin/providersdb examples/providers/cmd/providersdb/main.go
	GOGC=off go build -o bin/providersdb-gun examples/providers/cmd/providersdb-gun/main.go

generate:
	go generate


