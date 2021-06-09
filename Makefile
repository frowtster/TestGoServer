all: deps clean build

init:
	rm -f go.mod
	go mod init my-simple-server
	go mod tidy
	go mod verify

deps:
	go get -v github.com/julienschmidt/httprouter
	go get -v github.com/sirupsen/logrus

clean:
	rm -f server

build: deps
	go build -o server

run-dev:
	./server

.PHONY: \
	deps \
	clean \
	build \
	run-dev
