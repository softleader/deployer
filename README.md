# Deployer

Deployer is a tool for managing SoftLeader DevOps pipelines.

## Installation

```
$ go get -u github.com/kataras/iris
$ go get -u github.com/softleader/deployer
$ go run $GOPATH/src/github.com/softleader/deployer/main.go
```

## Usage

- List all stacks

```
$ curl localhost:5678
```

- List services in stack

```
$ curl localhost:5678/${stack}
```

- Remove stack ${stack}

```
$ curl -X DELETE localhost:5678/${stack}
```

- Deploy a `package.yaml`

```
$ curl -X POST \
       -d '{"project": "cki", "eurekaPort":8081, "kibanaPort":8000, "gatewayPort": 8080, "publishPort": 30000, "volume0": "/nfs/rpc", "net0": "softleader-cki", "yaml": "github:softleader/softleader-package/package.yaml#hotains"}' \
       localhost:5678/
```
