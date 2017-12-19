# Deployer

Deployer is a tool for managing SoftLeader DevOps pipelines.

## Installation

```
$ go get -u github.com/softleader/deployer
$ go build -o deployer/main github.com/softleader/deployer
$ ./deployer/main
```

### Args

- `wd` - Determine a working dictionary, default: `$(pwd)`
- `addr` - Determine application addr, default: `empty`
- `port` - Determine application port, default: `5678`
- `cmd.gpm` - Command to execute softleader/git-package-manager, default: `gpm`
- `cmd.-V
` - Command to execute softleader/container-yaml-generator, default: `gen-yaml`

eg.

```
$ ./main -wd=/tmp -port=8080
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
       -d '{"cleanUp": true, "project": "hotains", "eurekaPort":8081, "kibanaPort":8000, "gatewayPort": 8080, "dev": "192.168.1.60/30000", "volume0": "/nfs/rpc", "net0": "", "yaml": "package.yaml"}' \
       localhost:5678/
```
