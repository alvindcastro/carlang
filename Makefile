.PHONY: test build run-hello compile-hello vm-hello clean

test:
	go test ./...

build:
	go build -o bin/carl ./cmd/carl

run-hello:
	go run ./cmd/carl run examples/hello_world.carl

compile-hello:
	go run ./cmd/carl compile examples/hello_world.carl -o dist/hello_world.cbc

vm-hello: compile-hello
	go run ./cmd/carl vm dist/hello_world.cbc

clean:
	rm -rf bin dist
