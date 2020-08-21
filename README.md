# auto-build-go

由于需要经常新建 `Golang` 项目，打算写个脚手架，一键自动生成指定格式的项目，方便后续开发。   
完成项目中的所有的 `TODO`，你将得到一个能正常使用的项目。

### 项目结构

**注：项目可选项全部加上时结构**

```
project 
│
└─── cmd
│   └─── root.go      # 修改此处 TODO，添加需要执行的服务
│   └─── root.go
└─── config           # 配置文件文件夹
│   └─── config.go    # 修改 TODO，添加所需要的配置实体
│   └─── config.yaml  # 修改 config.go 中的配置实体后，需要同步添加默认值
└─── grpc             # GRPC 文件夹，在这个文件夹里面添加你的 proto，并实现。
│   └─── root.go
└─── pkg
│   └─── ***
└─── server           # 所有 cmd 所需要使用的服务应该在这里统一
│   └─── grpc.go
|   └─── proxy.go
└─── Makefile
```

### 使用方式

```shell
python3 main.py [-h] [--base BASE] [-o] [--go_path GO_PATH] [--grpc] createProjectName
# -h                帮助
# --base disable    禁止创建基础模块，不推荐
# -o --organization 项目组织名，默认：github.com
# --go_path         手动设置 gopath 地址，否则将查找环境变量，~/go/src 优先级最高
# --grpc            自动添加 grpc 模块
```
