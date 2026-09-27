.PHONY: run build test fmt clean

run:
	go run ./src

build:
	go build -o bin/emulator ./src

test:
	go test ./tests/...

fmt:
	gofmt -w src tests

clean:
	rm -rf bin
