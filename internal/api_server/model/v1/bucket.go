package v1

type Bucket struct {
	BucketID BucketID
	UserID   UserID
	Name     string
}

type BucketID uint64
