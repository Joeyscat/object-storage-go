
# 编译
FROM golang:1.16.3-alpine3.13 AS builder
#FROM golang:1.16.3-stretch AS builder
COPY . .
RUN go build -o storage_server ./cmd/storage_server/storage_server.go


# 构建运行程序的镜像
FROM alpine:3.12
COPY --from=builder /go/storage_server ./storage_server
#COPY index.html index.html

# TODO 配置环境变量
ENV RABBITMQ_SERVER=amqp://test:test@ppi.io:5672
ENV LISTEN_ADDRESS=localhost:8000
ENV STORAGE_ROOT=/tmp/object-storage/

# 运行
CMD ["./storage_server"]