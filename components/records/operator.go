package records

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/tags"
)

type ID common.IDStr

type Content struct {
	Title    string    `json:",omitempty" bson:",omitempty"`
	Summary  string    `json:",omitempty" bson:",omitempty"`
	TypeKey  string    `json:",omitempty" bson:",omitempty"`
	Data     string    `json:",omitempty" bson:",omitempty"`
	Embedded []Content `json:",omitempty" bson:",omitempty"` // in particular: URLs, images, etc.
}

type Item struct {
	ID                         ID      `json:",omitempty" bson:"_id,omitempty"`
	Content                    Content `json:",inline"    bson:",inline"`
	structures.ItemDescription `        json:",inline"    bson:",inline"`
}

type Operator interface {
	Save(Item, *auth.Identity) (ID, error)
	Read(ID, *auth.Identity) (*Item, error)
	Remove(ID, *auth.Identity) error

	List(*selectors.Term, *auth.Identity) ([]Item, error)
	Stat(*selectors.Term, *auth.Identity) (db.StatMap, error)
	Tags(*selectors.Term, *auth.Identity) (tags.StatMap, error)

	AddParent(ts []tags.Item, id ID) ([]tags.Item, error)
}
