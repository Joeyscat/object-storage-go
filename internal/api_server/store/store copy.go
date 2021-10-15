package store

type Store1 interface {
	// SearchLatestVersion 查询最新版本的元数据,如果元数据不存在,返回版本为0的元数据
	SearchLatestVersion(name string) (meta *ObjectMeta, err error)

	// GetByVersion 查询元数据,当version=0时,查询最新记录
	GetByVersion(name string, version int) (meta *ObjectMeta, err error)

	// SaveNewObject 插入元数据
	SaveNewObject(name, hash string, size uint64) (err error)

	// AddNewVersion 给元数据插入新版本记录
	AddNewVersion(name, hash string, version, size uint64) (err error)

	// SearchAllVersions 查询历史版本信息
	// 分页查询, from 从0开始
	SearchAllVersions(name string, from, size int64) (metas []*ObjectMeta, err error)

	// DelByVersion 根据对象名与版本号删除对象的元数据
	DelByVersion(name string, version int)

	// SearchVersionStatus 查询版本数量超过 minVersionCount 的元数据
	SearchVersionStatus(VersionCount int) ([]Bucket, error)

	// HasHash 查询元数据中是否存在该hash的对象
	HasHash(hash string) (bool, error)

	// SearchHashSize 获取哈希值对应的对象大小
	SearchHashSize(hash string) (size int64, err error)

	// SearchHash 根据hash查询对象
	SearchHash(hash string) (*ObjectMeta, error)
}

type ObjectMeta struct {
	Name    string // 对象的名字，不会改变
	Version uint32 // 对象当前的版本
	Size    uint64 // 对象当前的大小，单位字节
	Hash    string // 对象当前的哈希值
}

type Bucket struct {
	Name         string // 对象的名字
	VersionCount uint8  // 对象有多少个版本
	MinVersion   uint32 // 对象当前最小版本号
}
