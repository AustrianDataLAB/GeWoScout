build:
	go build -o bin/handler github.com/AustrianDataLAB/GeWoScout/backend

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/handler.exe github.com/AustrianDataLAB/GeWoScout/backend

all: build build-windows
