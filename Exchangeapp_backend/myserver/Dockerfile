# 使用官方golang:alpine镜像作为构建环境
FROM golang:alpine AS builder
#构建可执行文件
#关联CGO
#ENV CGO_ENABLED 0
#ENV GOPROXY https://goproxy.cn,direct
##设置yml的镜像源
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories


# 设置工作目录
RUN mkdir "/build"
WORKDIR /build
ADD go.mod .
ADD go.sum .
ADD main.go .
#运行命令打出包来,位置在上面设置的/build里面
RUN go build -o main


# 使用基本的scratch镜像作为最终镜像 启动一个容器scratch非常小 工作目录在/app下
FROM scratch
WORKDIR /app
COPY ./config.yaml .
# 将编译好的二进制文件从构建阶段复制到最终镜像中（把/build/main文件copy到/app下面）
COPY --from=builder /build/main /app
#设置容器启动时执行的命令
CMD ["./main"]
