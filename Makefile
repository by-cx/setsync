all: build

build:
	dep ensure
	go build -o setsync *.go