package store

var client Factory

// Factory defines the iam platform storage interface.
type Factory interface {
	Buckets() BucketStore
	Objects() ObjectStore
	Close() error
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the store client.
func SetClient(factory Factory) {
	client = factory
}
