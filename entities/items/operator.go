package items

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/files"
)

type Items = files.Items
type Item = files.Item

type Operator interface {
	Save(path, newFilePattern string, data []byte, identity *auth.Identity) (string, error)
	Read(path string, identity *auth.Identity) ([]byte, error)
	Remove(path string, identity *auth.Identity) error
	List(path string, depth int, identity *auth.Identity) (Items, error)
	Stat(path string, depth int, identity *auth.Identity) (*Item, error)
}
