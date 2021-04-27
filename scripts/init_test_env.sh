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

rm -r /tmp/object-storage/

for i in $(seq 1 6)
do
  mkdir -p /tmp/object-storage/$i/objects
  mkdir -p /tmp/object-storage/$i/temp
  mkdir -p /tmp/object-storage/$i/garbage
done
