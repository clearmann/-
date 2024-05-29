
# 使用官方Go镜像作为构建环境
FROM golang:1.20-alpine as builder

# 设置工作目录
WORKDIR /app

# 拷贝go.mod和go.sum文件
COPY go.mod go.sum ./
COPY config.yaml ./
# 下载依赖
RUN go mod download

# 拷贝源代码文件
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# 继续使用Alpine镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段拷贝编译好的可执行文件
COPY --from=builder /app/main /app/main

# 暴露应用监听的端口号
EXPOSE 8080
# 运行程序
CMD ["ENV export GIN_MODE=release"]
CMD ["./main"]