GOPATH=$(CURDIR)

.PHONY: default prepare get-deps

default:
	env GOPATH=$(CURDIR) go build -o bin/gmanga gmanga

get-deps:
	env GOPATH=$(CURDIR) go get -u -v github.com/gorilla/mux
	env GOPATH=$(CURDIR) go get -u -v github.com/gorilla/handlers
