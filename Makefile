all: build 

build: 
	mkdir -p $(GOPATH)
	go get -v github.com/ncode/hot-potato-router/...
	go build -v ./...
	go build -o $(GOPATH)/bin/hot-potato-router .
	
install:
	mkdir -p $(DESTDIR)/usr/bin/
	mkdir -p $(DESTDIR)/etc/hpr/
	cp -a $(GOPATH)/bin/hot-potato-router $(DESTDIR)/usr/bin/hpr
	cp config/hpr.yml $(DESTDIR)/etc/hpr/
