.PHONY: all clean generate build

all: clean generate build

clean:
	rm -rf bin

build:
	GOGC=off go build -o bin/dronesdb examples/drones/cmd/dronesdb/main.go
	GOGC=off go build -o bin/dronesdb-bench examples/drones/cmd/dronesdb-bench/main.go

generate:
	go generate


