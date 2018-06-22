GOARCH=amd64
BINARY=build

all: build npm docker clean

build:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}/main .

npm:
	npm install

docker:
	docker build -t softleader/deployer .

publish:
	docker push softleader/deployer

clean:
	rm -rf ${BINARY} node_modules/ package-lock.json

.PHONY: clean build