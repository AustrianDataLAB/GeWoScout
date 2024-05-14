build:
	go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/handler.exe github.com/AustrianDataLAB/GeWoScout/backend

build-linux-minimal:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

build-deployment:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -buildid=" -trimpath -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

run: build
	func start

swagger:
	swag init -g handler.go

all: build build-windows
