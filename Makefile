all: build

build:
	GOPATH=$(GOPATH) go get -v github.com/ncode/hot-potato-router/...
	GOPATH=$(GOPATH) go build -o $(GOPATH)/bin/hpr $(GOPATH)/src/github.com/ncode/hot-potato-router/hpr.go
	GOPATH=$(GOPATH) go build -o $(GOPATH)/bin/hprctl $(GOPATH)/src/github.com/ncode/hot-potato-router/hprctl/hprctl.go

install:
	mkdir -p $(DESTDIR)/usr/bin/
	mkdir -p $(DESTDIR)/etc/hpr/
	cp -a $(GOPATH)/bin/hpr $(DESTDIR)/usr/bin/hpr
	cp -a $(GOPATH)/bin/hprctl $(DESTDIR)/usr/bin/hprctl
	cp config/hpr.yml $(DESTDIR)/etc/hpr/
