.PHONY: swagger

build:
	go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/handler.exe github.com/AustrianDataLAB/GeWoScout/backend

run: build
	func start

swagger:
	swag init -g handler.go

build-ci: swagger build

all: build build-windows
