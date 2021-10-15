# start up processes for test

# export MONGODB_URI=mongodb://object_storage_rw:123456@192.168.50.186:27017/object_storage
export MONGODB_URI=mongodb://object_storage_rw:5QXj_hVQ7_5r5oOr1KVXjGam00qVgCZ35d5BmTxTYDpemN4d7o7SxCp1euiGtCR3@mongo1.me.io:20001,mongo2.me.io:20001,mongo3.me.io:20001/object_storage
export NATS_URL=nats://nats.me.io:4222
export NATS_SUBJECT_STORAG_HEARTBEAT=storage_heartbeat
export NATS_SUBJECT_OBJ_LOCATE=object_locate

BASE_DIR=~/env/objects
LOG_DIR=~/code/go/object-storage-go/logs

LISTEN_ADDRESS=localhost:8001 STORAGE_ROOT=$BASE_DIR/1 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-1.log 2>&1 &
LISTEN_ADDRESS=localhost:8002 STORAGE_ROOT=$BASE_DIR/2 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-2.log 2>&1 &
LISTEN_ADDRESS=localhost:8003 STORAGE_ROOT=$BASE_DIR/3 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-3.log 2>&1 &
LISTEN_ADDRESS=localhost:8004 STORAGE_ROOT=$BASE_DIR/4 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-4.log 2>&1 &
LISTEN_ADDRESS=localhost:8005 STORAGE_ROOT=$BASE_DIR/5 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-5.log 2>&1 &
LISTEN_ADDRESS=localhost:8006 STORAGE_ROOT=$BASE_DIR/6 go run cmd/storage_server/storage_server.go >> $LOG_DIR/data-server-6.log 2>&1 &


LISTEN_ADDRESS=localhost:9001 go run cmd/api_server/api_server.go >> $LOG_DIR/api-server-1.log 2>&1 &
LISTEN_ADDRESS=localhost:9002 go run cmd/api_server/api_server.go >> $LOG_DIR/api-server-2.log 2>&1 &

# =============================================================================
echo 'waiting ...'

sleep 5s
# request api

#echo -n "this object will have only 1 instance" | openssl dgst -sha256 -binary | base64
#aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=

curl -v localhost:9001/objects/test4_1 -XPUT -d"this object will have only 1 instance" -H "Digest: SHA-256=aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="

curl localhost:9001/locate/test4_1

echo
curl localhost:9002/objects/test
echo

# =============================================================================

# stop processes

# kill $(ps -ef | grep go-build | awk '$0 !~/grep/ {print $2}' | tr -s '\n' ' ')
# kill $(ps -ef | grep server.go | awk '$0 !~/grep/ {print $2}' | tr -s '\n' ' ')


