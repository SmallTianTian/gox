FROM golang:1.17 AS build

# 指定工作目录
WORKDIR /build

# 拷贝整个项目的文件
COPY .  .

# 提前做缓存层
RUN go mod download

# 使用外部编译变量
ARG version

# 定义匿名卷，复用 cache
# 可修改挂载点，和宿主机保持 cache 一致
VOLUME /go

# 编译，和 Makefile 中保持一致
RUN go build -a -installsuffix cgo -o gox -ldflags \
    "-s -w -X 'main.Version=${version}' \
    -X 'main.BuildTime=`date "+%Y-%m-%d %H:%M:%S"`' \
    -X 'main.GoVersion=`go version`' \
    -X 'main.GitRevision=`git rev-parse HEAD`' \
    -extldflags -static" .

# 二阶段编译
FROM golang:1.17 AS prod

# 从上一阶段中拷贝出编译好的文件
COPY --from=build /build/gox /bin/gox

# 指定工作目录
WORKDIR /workspace

# 设置 OS 为 docker
ENV OS=docker

# 默认命令
CMD ["gox", "help"]
