package catalogue_www

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/pavlo67/tools/entities/files"
)

type Item struct {
	URN  ns.URN
	File files.Item
}

type Operator interface {
	Save(bucketID files.BucketID, item Item, identity *auth.Identity) (string, error)
	Read(bucketID files.BucketID, path string, identity *auth.Identity) (*Item, error)
	Remove(bucketID files.BucketID, path string, identity *auth.Identity) error
	List(bucketID files.BucketID, path string, depth int, identity *auth.Identity) ([]Item, error)
	Stat(bucketID files.BucketID, path string, depth int, identity *auth.Identity) (*Item, error)
}
