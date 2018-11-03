all: build

dep:
	dep ensure

build: dep test
	go build -o setsync *.go

test: dep
	go test