package catalogue

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"

	"github.com/pavlo67/tools/components/files"
	"github.com/pavlo67/tools/components/records"
)

type ID common.IDStr

type Item struct {
	ID     ID
	Record records.Item
	File   files.Item
}

type Operator interface {
	Save(bucketID files.BucketID, item Item, options *crud.Options) (*ID, error)
	Read(bucketID files.BucketID, id ID, options *crud.Options) (*Item, error)
	Remove(bucketID files.BucketID, id ID, options *crud.Options) error
	List(bucketID files.BucketID, path string, depth int, options *crud.Options) ([]Item, error)
	Stat(bucketID files.BucketID, path string, depth int, options *crud.Options) (Item, error)
}
