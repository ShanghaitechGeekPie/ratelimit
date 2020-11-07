FROM golang:alpine

Maintainer  Leomund & Canarypwn @GeekPie_Association

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件app
RUN go build

# 启动容器时运行的命令
CMD ["/build/ratelimit"]
