
# 编译
FROM golang:1.16.3-alpine3.13 AS builder
#FROM golang:1.16.3-stretch AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o storage-server ./cmd/storage_server/storage_server.go && CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/api_server/api_server.go

# 构建运行程序的镜像
FROM alpine:3.12
COPY --from=builder /app/storage-server /app/api-server /

# 存储服务需要
RUN mkdir data data/objects data/temp data/garbage

# 配置环境变量
ENV RABBITMQ_SERVER=amqp://test:test@172.17.0.2:5672 LISTEN_ADDRESS=localhost:8000 STORAGE_ROOT=/data
# 暴露端口
EXPOSE 8000
# 默认运行storage-server
CMD ["/storage-server"]

# docker run --name os2 object-storage /api-server # Use this to startup api-server