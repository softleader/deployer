# Deployer

Deployer is a tool for managing SoftLeader DevOps pipelines.

## Installation

```
$ go get -u github.com/softleader/deployer
$ go build -o deployer/main github.com/softleader/deployer
$ ./deployer/main
```

### Args

- `workspace` - Determine a workspace, default: `$(pwd)/workspace`
- `addr` - Determine application addr, default: `empty`
- `port` - Determine application port, default: `5678`
- `cmd.gpm` - Command to execute softleader/git-package-manager, default: `gpm`
- `cmd.-V
` - Command to execute softleader/container-yaml-generator, default: `gen-yaml`

eg.

```
$ ./main -workspace=/tmp -port=8080
```

### Install as Ubuntu service

- Copy `deployer.service` to the directory `/etc/systemd/system/`
- Modify `ExecStart` in `deployer.service`
- 執行 `systemctl daemon-reload`
- Then it should be possible to control daemon using:

```
# 服務狀態
$ systemctl status deployer

# 服務開關
$ systemctl start deployer
$ systemctl stop deployer
$ systemctl restart deployer

# 開機自動啟動服務
$ systemctl enable deployer
$ systemctl disable deployer
```