#!/usr/bin/env bash

#for i in `seq 1 6`; do mkdir -p /tmp/$i/objects; done
export RABBITMQ_SERVER=amqp://test:test@ppi.io:5672
export ES_SERVER=localhost:9200

LISTEN_ADDRESS=localhost:8001 STORAGE_ROOT=/tmp/object-storage/1 go run chapter04/storage_server/cmd/server.go >> logs/data-server-1.log 2>&1 &
LISTEN_ADDRESS=localhost:8002 STORAGE_ROOT=/tmp/object-storage/2 go run chapter04/storage_server/cmd/server.go >> logs/data-server-2.log 2>&1 &
LISTEN_ADDRESS=localhost:8003 STORAGE_ROOT=/tmp/object-storage/3 go run chapter04/storage_server/cmd/server.go >> logs/data-server-3.log 2>&1 &
LISTEN_ADDRESS=localhost:8004 STORAGE_ROOT=/tmp/object-storage/4 go run chapter04/storage_server/cmd/server.go >> logs/data-server-4.log 2>&1 &
LISTEN_ADDRESS=localhost:8005 STORAGE_ROOT=/tmp/object-storage/5 go run chapter04/storage_server/cmd/server.go >> logs/data-server-5.log 2>&1 &
LISTEN_ADDRESS=localhost:8006 STORAGE_ROOT=/tmp/object-storage/6 go run chapter04/storage_server/cmd/server.go >> logs/data-server-6.log 2>&1 &


LISTEN_ADDRESS=localhost:9001 go run api_server/cmd/server.go >> logs/api-server-1.log 2>&1 &
LISTEN_ADDRESS=localhost:9002 go run api_server/cmd/server.go >> logs/api-server-2.log 2>&1 &
