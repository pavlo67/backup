package catalogue

import (
	"github.com/pavlo67/common/common/crud"

	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/pavlo67/tools/components/files"
)

type Item struct {
	IssuedID ns.URN
	File     files.Item
}

type Operator interface {
	Save(bucketID files.BucketID, item Item, options *crud.Options) (string, error)
	Read(bucketID files.BucketID, path string, options *crud.Options) (*Item, error)
	Remove(bucketID files.BucketID, path string, options *crud.Options) error
	List(bucketID files.BucketID, path string, depth int, options *crud.Options) ([]Item, error)
	Stat(bucketID files.BucketID, path string, depth int, options *crud.Options) (*Item, error)
}
