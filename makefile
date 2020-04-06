all: test build build-windows

build:
	go build

build-windows:
	cd cmd/relayserver && GOOS=windows GOARCH=amd64 go build

test:
	go test
