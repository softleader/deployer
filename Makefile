GOARCH = amd64

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o main-linux-${GOARCH} . ;

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o main-windows-${GOARCH}.exe . ;