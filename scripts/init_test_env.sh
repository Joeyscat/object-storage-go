## nats
# docker run -d -p 4222:4222 nats

## mongodb

## data and log dirs

BASE_DIR=~/env/objects
LOG_DIR=~/code/go/object-storage-go/logs

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
