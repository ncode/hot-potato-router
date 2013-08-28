all: install

install:
	go get -v github.com/ncode/hot-potato-router/...
	go build -o bin/hpr hpr.go
