package rs

const (
    DataShards    = 1
    ParityShards  = 1
    AllShards     = DataShards + ParityShards
    BlockPerShard = 8000
    BlockSize     = BlockPerShard * DataShards
)
