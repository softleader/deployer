GOARCH=amd64
BINARY=build
TAG=latest

all: build npm docker

build:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}/main .

npm:
	npm install -f

docker:
	docker build -t softleader/deployer:${TAG} .

publish:
	docker push softleader/deployer:${TAG}

clean:
	rm -rf ${BINARY} node_modules/ package-lock.json

.PHONY: clean build
