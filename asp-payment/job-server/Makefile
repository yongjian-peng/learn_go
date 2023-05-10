.PHONY: all build run gotool install clean help

all: gotool buildApi

buildApi:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server -ldflags="-w -s"

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make buildApi

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

