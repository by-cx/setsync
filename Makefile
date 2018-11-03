VERSION=0.1

.PHONY: all
all: build

.PHONY: dep
dep:
	dep ensure

.PHONY: build
build: dep test
	go build -o setsync *.go

.PHONY: release
release: dep build
	GOOS=linux GOARCH=amd64 go build -o setsync-${VERSION}.linux.amd64 -ldflags '-s' *.go
	GOOS=linux GOARCH=386 go build -o setsync-${VERSION}.linux.386 -ldflags '-s' *.go
	GOOS=linux GOARCH=arm go build -o setsync-${VERSION}.linux.arm -ldflags '-s' *.go
	
.PHONY: test
test: dep
	go test