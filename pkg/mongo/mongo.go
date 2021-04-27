package mongo

import (
    "context"
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/qiniu/qmgo"
    "github.com/qiniu/qmgo/middleware"
    "github.com/qiniu/qmgo/operator"
    "github.com/qiniu/qmgo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

var cli *qmgo.QmgoClient

func init() {
    InitCli()
}

func InitCli() {
    ctx := context.Background()
    var err error
    //cli, err = qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://db_object_storage_rw:jukBtxREsQo75AXg0gSnP7kbDh0fzQrN98Jq0mcAhwJAt0SXusPRtqAAlsb2WBir@10.176.93.248:20001/db_object_storage", Database: "db_object_storage", Coll: "t_metadata"})
    cli, err = qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://object_storage_rw:123456@localhost:27017/db_object_storage", Database: "db_object_storage", Coll: "t_metadata"})
    if err != nil {
        panic(err)
    }
    err = cli.Ping(2)
    if err != nil {
        panic(err)
    }

    // 创建索引
    so := &StorageObject{}
    err = cli.Find(context.Background(), bson.M{}).One(so)
    if err != nil && err == mongo.ErrNoDocuments {
        err = cli.
            CreateIndexes(context.Background(),
                []options.IndexModel{
                    {Key: []string{"name"}, Unique: true},
                    {Key: []string{"name", "versions.v"}, Unique: true},
                })
        if err != nil {
            panic(err)
        }
    }

    middleware.Register(func(doc interface{}, opType operator.OpType, opts ...interface{}) error {
        log.Info(fmt.Sprintf("doc: %v\n opType: %v\n opts: %v", doc, opType, opts))
        return nil
    })
}

type Metadata struct {
    Name    string
    Version uint64
    Size    uint64
    Hash    string
}

// StorageObject 存储对象在mongo中的文档结构
type StorageObject struct {
    Name     string     `bson:"name"`
    Versions []*Version `bson:"versions"`
}

// Version 存储对象版本信息
type Version struct {
    V    uint64 `bson:"v"`
    Size uint64 `bson:"size"`
    Hash string `bson:"hash"`
}

// getMetadata 获取元数据,返回的版本信息需要匹配指定的版本号.
// db.t_metadata.findOne({"name":"xxx1", "versions.v":2}, {"versions.$":1, "name":1})
func getMetadata(name string, versionId int) (meta *Metadata, err error) {
    so := new(StorageObject)
    if err = cli.Find(context.Background(), bson.M{"name": name, "versions.v": versionId}).
        Select(bson.M{"versions.$": 1, "name": 1}).One(so); err == mongo.ErrNoDocuments {
        return &Metadata{ // 返回无效的元数据
            Name:    name,
            Version: 0,
            Size:    0,
            Hash:    "",
        }, nil
    }
    if err != nil {
        return nil, err
    }
    if so.Versions != nil && len(so.Versions) > 0 {
        return &Metadata{
            Name:    name,
            Version: so.Versions[0].V,
            Size:    so.Versions[0].Size,
            Hash:    so.Versions[0].Hash,
        }, nil
    }

    return nil, fmt.Errorf("查询不到[%s]元数据的版本信息", name)
}

// SearchLatestVersion 查询最新版本的元数据,如果元数据不存在,返回版本为0的元数据
//
// db.t_metadata.aggregate([
//     {$match: {"name":"xxx1"}},
//     {$project: {"versions": 1, "name": 1}},
//     {$unwind: "$versions"},
//     {$sort: {"versions.v": -1}},
//     {$limit: 1}
// ])
func SearchLatestVersion(name string) (meta *Metadata, err error) {
    type LatestVersion struct {
        Name     string   `bson:"name"`
        Versions *Version `bson:"versions"`
    }
    l := new(LatestVersion)

    err = cli.Aggregate(context.Background(), mongo.Pipeline{
        bson.D{{"$match", bson.M{"name": name}}},
        bson.D{{"$project", bson.M{"versions": 1, "name": 1}}},
        bson.D{{"$unwind", "$versions"}},
        bson.D{{"$sort", bson.M{"versions.v": -1}}},
        bson.D{{"$limit", 1}},
    }).One(l)

    if err != nil && err == mongo.ErrNoDocuments {
        return &Metadata{ // 返回无效的元数据
            Name:    name,
            Version: 0,
            Size:    0,
            Hash:    "",
        }, nil
    }

    if err != nil {
        return nil, err
    }
    if l.Versions != nil {
        meta = &Metadata{
            Name:    name,
            Version: l.Versions.V,
            Size:    l.Versions.Size,
            Hash:    l.Versions.Hash,
        }
        return meta, nil
    }

    return nil, fmt.Errorf("[%s]的元数据没有版本信息", name)
}

// GetMetadata 查询元数据,当version=0时,查询最新记录
func GetMetadata(name string, version int) (meta *Metadata, err error) {
    if version == 0 {
        return SearchLatestVersion(name)
    }
    return getMetadata(name, version)
}

// PutMetadata 插入元数据
func PutMetadata(name, hash string, size uint64) (err error) {
    so := &StorageObject{
        Name: name,
        Versions: []*Version{{
            V:    1,
            Size: size,
            Hash: hash,
        }},
    }
    result, err := cli.InsertOne(context.Background(), so)
    if err != nil {
        return err
    }
    log.Debug(fmt.Sprintf("%v", result))

    return
}

// AddVersion 给元数据插入新版本记录
func AddVersion(name, hash string, version, size uint64) (err error) {
    coll, err := cli.Collection.CloneCollection()
    if err != nil {
        return err
    }

    filter := bson.M{"name": name, "versions.v": bson.M{"$ne": version}}
    update := bson.M{"$push": bson.M{"versions": Version{
        V: version, Size: size, Hash: hash,
    }}}
    res, err := coll.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return fmt.Errorf("找不到[%s]的元数据或已存在该版本[%d],请将v+1并重试", name, version)
    }
    return
}

// SearchAllVersions 查询历史版本信息
// 分页查询, from 从0开始
func SearchAllVersions(name string, from, size int64) (metas []*Metadata, err error) {
    so := new(StorageObject)

    query := bson.M{"name": name}
    projection := bson.M{"versions": bson.M{"$slice": bson.A{from, size}}}
    err = cli.Find(context.Background(), query).Select(projection).One(so)
    if err != nil && err == mongo.ErrNoDocuments {
        return nil, nil
    }

    for _, version := range so.Versions {
        metas = append(metas, &Metadata{
            Name:    name,
            Version: version.V,
            Size:    version.Size,
            Hash:    version.Hash,
        })
    }

    return
}

// DelMetadata 根据对象名与版本号删除对象的元数据
func DelMetadata(name string, version int) {
    panic("unimplemented")
}

type Bucket struct {
    Key        string
    DocCount   int
    MinVersion struct {
        Value float32
    }
}

type aggregateResult struct {
    Aggregations struct {
        GroupByName struct {
            Buckets []Bucket
        }
    }
}

// SearchVersionStatus 查询版本数量超过 minDocCount 的元数据
// 返回的数据结构包含对象的名字,该对象有多少个版本,该对象当前最小版本号
func SearchVersionStatus(minDocCount int) ([]Bucket, error) {
    panic("unimplemented")
}

// HasHash 查询元数据中是否存在该hash的对象
func HasHash(hash string) (bool, error) {
    panic("unimplemented")
}

// SearchHashSize 获取哈希值对应的对象大小
func SearchHashSize(hash string) (size int64, err error) {
    panic("unimplemented")
}

// SearchHash 根据hash查询对象
// db.t_metadata.find({"versions.hash":"xxx"},{"versions":{$slice:[0,1]}})
func SearchHash(hash string) (*Metadata, error) {

    return nil, nil
}
