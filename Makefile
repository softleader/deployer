GOARCH = amd64

all: linux windows

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o main-linux-${GOARCH} . ;

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o main-windows-${GOARCH}.exe . ;