# fresh-go

> 新建一个项目，我们需要很多非业务代码相关的文件，它们可能来自其他项目，也可能来自自己的想法，但要正常运行起来，往往需要花费我们很多时间，对提高工作效率毫无帮助。

fresh-go 是一个快速项目生成器，在设计上，它遵循如下准则：
- 简单。生成完项目即可运行，快速查看效果。
- 高效。通过简单的命令即可完成诸如 gRPC proto 创建等需求。
- 通用。基础库不自建，使用 gin、ent、viper、gRPC 等知名框架，没有换框架的负担。
fresh-go 不仅仅只是一个快速项目生成器，它也将引入 git commit 规范、热调试等工具，立志成为你开发路上的良师益友。

文档地址：https://fresh-go.tianxu.xin

**小贴士**：
1. go verion >= 1.6
2. grpc 编译工具将使用 v2 版本


## 简单开始：
### 安装本工具
```shell
$ go get -u github.com/SmallTianTian/fresh-go@master
```

### 安装本工具依赖 go 相关的插件
**注意：执行该命令后，grpc 等编译工具将被替换为 v2 版本，可能导致历史项目编译变动很大。**
```shell
$ fresh-go init
```

### 创建项目
创建一个 [module](https://blog.golang.org/using-go-modules) 为 `github.com/SmallTianTian/helloworld` 的项目
```shell
$ fresh-go new helloworld \          # 项目名
             --remote=github.com \   # remote 地址
             --owner=SmallTianTian   # owner 名称
```

### 运行 demo
```shell
$ make run # 推荐做法，但需要本地有 make 及 golang-ci 插件
# 上下两条命令选一条执行即可
$ go run cmd/server/main.go cmd/server/wire_gen.go # 仅使用 go 命令
```
正常运行项目后，可以通过如下命令查看效果。
```shell
$ curl 'http://localhost:8089/say/bob'

# 下面为响应示例
# {"message":"Hello bob~"}
```

### 新建 API
当需要加入新的服务时，你可以通过本工具来快速创建模板
```shell
$ fresh-go new user
```
添加完毕后，可以刚刚运行 demo 的命令运行程序，在命令行中通过如下命令查看效果。
```shell
$ curl -v 'http://localhost:8089/v1/user/1'

# 下面为新建 API，但没有实现的响应
# {"Code":1,"Message":"NOT SUPPORT.","Action":"","Reason":"","Details":[],"CauseError":""}
```
后续补齐 proto 中的字段，完善 `internal/ui/grpc/user_v1.go` 中的服务即可再次查看效果。
