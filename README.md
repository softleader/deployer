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
curl localhost:5678
```

- Remove stack ${stackName}

```
curl -X DELETE localhost:5678/${stackName}
```

- Deploy a `package.yaml`

```
curl -X POST -d '{"eurekaPort":8081, "kibanaPort":8000, "gatewayPort": 8080, "publishPort": 30000, "mount": "/nfs/rpc", "network": "net0", "yaml": "github:softleader/softleader-package/package.yml#rpc"}' localhost:5678/${stackName}
```
