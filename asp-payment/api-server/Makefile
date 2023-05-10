.PHONY: all build run gotool install clean help

all: gotool buildApi

buildApi:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server -ldflags="-w -s"
deployApi:
	cd /data/sunnypay/common && git pull
	cd /data/sunnypay/api-server && git pull
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o payment-api
	rm -rf /data/payment-api/payment-api
	mv payment-api /data/payment-api/
	systemctl restart payment-api
	echo "deploy Api success"
run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make buildApi

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

