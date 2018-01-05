GOARCH=amd64
BINARY=build

all: clean linux windows

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}/main-linux-${GOARCH} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}/main-windows-${GOARCH}.exe .

clean:
	rm -rf ${BINARY}

.PHONY: clean