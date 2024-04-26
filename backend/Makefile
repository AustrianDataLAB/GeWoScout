ifeq ($(OS),Windows_NT)
	TF_BUILD_CMD := SET GOOS=linux&&SET GOARCH=amd64&&go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend
else
	TF_BUILD_CMD := GOOS=linux GOARCH=amd64 go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend
endif

build:
	go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/handler.exe github.com/AustrianDataLAB/GeWoScout/backend

run: build
	func start

swagger:
	swag init -g handler.go

all: build build-windows

build-tf:
	${TF_BUILD_CMD}