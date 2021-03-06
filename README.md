# Deployer

> Deployer is a tool for managing SoftLeader DevOps pipelines.

## Installation

```shell
$ go get -u github.com/softleader/deployer
```

### Run

```shell
$ cd $GOPATH/src/github.com/softleader/deployer
$ make
$ ./build/main-macos-amd64 # main-linux-amd64, main-windows-amd64.exe
```

open [http://localhost:5678](http://localhost:5678)

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

### Docker Compose

- Write a docker-compose YAML file

```yml
deployer:
  image: "softleader/deployer"
  ports:
    - 5678:80
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
    - ./.gpm:/root/.gpm
    - ./workspace:/workspace
  restart: always
```

> `docker-compose.yml` stores in `/devops/deployer` on 192.168.1.60

- Run via docker compose

```shell
docker-compose up -d
```

### [Deprecated] Install as Ubuntu service (on 192.168.1.60)

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

### [Deprecated] Update deployer service (on 192.168.1.60)

```shell
$ sudo su
$ go get -u github.com/softleader/deployer
$ cd /root/go/src/github.com/softleader/deployer/
$ make linux
$ systemctl restart deployer
```

## REST API

> 讓 Deployer 可以快速的串接在其他 CI/CD 系統的 pipeline 中

### Deploy your package.yaml

```
POST /api/stacks
```

#### Parameter

| Name | Type | Description | Default Value |
|------|------|-------------|---------------|
| `cleanUp` | `bool` | 執行前是否要先清空 workspace | false |
| `project` | `string` | **Required.** 用來組 stack name, 組成方式為: `${project}-${port}-${group}` ( `-${port}` 及 `-${group}` 視部署條件不同決定是否會加上) | |
| `yaml` | `string` | **Required.** yaml 位置 | |
| `style` | `string` | 部署 yaml 的樣式: swarm, k8s | k8s |
| `volume0` | `string` | `Containerfile` 中 `volumn0` 變數的實際目錄 | |
| `net0` | `string` | `Containerfile` 中 `net0` 的 network 名稱, 無傳值則由 docker 自動建立 network | |
| `group` | `string` | 指定要部署的 group, 可用`,`串起多個, 沒傳入則部署所有 group | |
| `flatGroup` | `bool` | 是否要打平 group, 讓所有 group 都部署在同一個 stack 中 | false |
| `silently` | `bool` | 安靜模式 | false |
| `dev.ipAddress` | `string` | dev 模式的 ip, 有傳值則開啟 dev 模式, 反之則關閉 dev 模式 | |
| `dev.port` | `int` | dev 模式下的 port | |
| `dev.ignore` | `string` | dev 模式下要忽略的 group 名稱, 可用`,`串起多個 | |

```json
{
    "cleanUp": true,
    "project": "example",
    "yaml": "github:softleader/softleader-package/softleader-base.yaml#master",
    "style": "swarm",
    "net0": "example-network",
    "volume0": "",
    "group": "",
    "flatGroup": false,
    "silently": false,
    "dev": {
    	"ipAddress": "192.168.1.60",
    	"port": 50001,
    	"ignore": "elasticsearch,kibana,logstash"
    }
}
```

##### Example

```sh
$ curl -X POST \
    -d '{
      "cleanUp":true,
      "project":"example",
      "style":"swarm",
      "yaml":"github:softleader/softleader-package/softleader-base.yaml#master"
    }' \
    http://softleader.com.tw:5678/api/stacks
```

### Remove a Stack

```
DELETE /api/stacks/:stack
```

`:stack` - 左右模糊比對, 刪除所有符合的 stack 名稱

##### Example

```
$ curl -X DELETE http://softleader.com.tw:5678/api/stacks/:stack
```

### Remove a Service

```
DELETE /api/services/:service
```

`:service` - 完整比對, 可以是 service id 或 service name

##### Example

```
$ curl -X DELETE http://softleader.com.tw:5678/api/services/:service
```

### Update Service image or replicas


```
PUT /api/services/:service
```

`:service` - 完整比對, 可以是 service id 或 service name, 或是 `filter` 的條件

#### Parameter

| Name | Type | Description | Constraint |
|------|------|-------------|---------------|
| `image` | `string` | 要更新的 image | `image` 及 `replicas` 至少必須給其一 |
| `replicas` | `int` | 要更新的 replicas | `image` 及 `replicas` 至少必須給其一 |
| `filter` | `string` | 使用傳入的 filter 配上 `:service`, 來過濾出要跟新的 service id | filter 的條件必須要可以過濾出唯一的 service |
|  `skip-slack` | `any` | 過版時不要 hook slack | | 

##### Example

```sh
# 將 id 為 :service 的更新成 busybox:1.28 並開 2 個 replicas
$ curl -X PUT 'http://softleader.com.tw:5678/api/services/:service?image=busybox:1.28&replicas=2'

# 將 service label 為 app=busybox 的 service 更新成 1 個 replicas, 並且不要通知 slack
$ curl -X PUT 'http://softleader.com.tw:5678/api/services/app=busybox?filter=label&replicas=1&skip-slack'
```

