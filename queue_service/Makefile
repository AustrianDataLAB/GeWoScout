build:
	go build -o bin/handler handler.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/handler.exe handler.go

all: build build-windows
