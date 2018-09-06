.PHONY: all deps test lint

all: deps lint test

check-gometalinter:
	which gometalinter || (go get -u -v github.com/alecthomas/gometalinter && gometalinter --install)

deps:
	go get -t -v ./...

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

lint: check-gometalinter
	gometalinter --exclude=vendor --disable-all --enable=golint --enable=vet --enable=vetshadow --enable=gofmt ./...
