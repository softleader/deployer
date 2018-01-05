# Deployer

Deployer is a tool for managing SoftLeader DevOps pipelines.

## Installation

```shell
$ go get -u github.com/softleader/deployer
$ cd $GOPATH/src/github.com/softleader/deployer
$ make
$ ./build/main-macos-amd64 # main-linux-amd64, main-windows-amd64.exe
```

### Args

- `workspace` - Determine a workspace, default: `$(pwd)/workspace`
- `addr` - Determine application addr, default: `empty`
- `port` - Determine application port, default: `5678`
- `cmd.gpm` - Command to execute [softleader/git-package-manager](https://github.com/softleader/git-package-manager), default: `gpm`
- `cmd.gen-yaml` - Command to execute [softleader/container-yaml-generator](https://github.com/softleader/container-yaml-generator), default: `gen-yaml`

eg.

```shell
$ ./build/main-macos-amd64 -workspace=/tmp -port=8080
```

### Install as Ubuntu service (192.168.1.60)

- Copy `deployer.service` to the directory `/etc/systemd/system/`

  ````shell
  $ cp /devops/deployer/deployer.service /etc/systemd/system/
  ````

- Reload systemd manager configuration

  ````shell
  $ systemctl daemon-reload
  ````

- Then it should be possible to control daemon using:

  ```shell
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

### Update deployer service (192.168.1.60)

```shell
$ sudo su
$ go get github.com/softleader/deployer
$ systemctl restart deployer
```
