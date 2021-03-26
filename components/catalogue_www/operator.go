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
	Save(item Item, identity *auth.Identity) (string, error)
	Read(path string, identity *auth.Identity) (*Item, error)
	Remove(path string, identity *auth.Identity) error
	List(path string, depth int, identity *auth.Identity) ([]Item, error)
	Stat(path string, depth int, identity *auth.Identity) (*Item, error)
}
