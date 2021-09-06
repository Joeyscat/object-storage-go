#wget ppi.io:15672/cli/rabbitmqadmin
#chmod +x rabbitmqadmin
#mv rabbitmqadmin /usr/bin/
#
#rabbitmqadmin --host=ppi.io declare exchange name=api-server type=fanout
#rabbitmqadmin --host=ppi.io declare exchange name=data-server type=fanout
#rabbitmqadmin --host=ppi.io list exchanges
#
## docker pull elasticsearch:7.10.1
#docker run -d --name es-object-storage -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch
#curl ppi.io:9200/metadata -H "Content-Type: application/json" -XPUT -d'{"mappings":{"properties":{"name":{"type":"text"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"text"}}}}'

BASE_DIR=/tmp/object-storage
LOG_DIR=/tmp/object-storage/logs

if [ ! -d "$BASE_DIR" ]; then
    echo "x"
else
    rm -r $BASE_DIR
fi

if [ ! -d "$LOG_DIR" ]; then
    mkdir -p $LOG_DIR
fi


for i in $(seq 1 6)
do
  mkdir -p $BASE_DIR/$i/objects
  mkdir -p $BASE_DIR/$i/temp
  mkdir -p $BASE_DIR/$i/garbage
done
