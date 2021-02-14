package records

import (
	"time"

	"github.com/pavlo67/tools/components/tags"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/tools/components/ns"
	"github.com/pavlo67/tools/components/vcs"
)

type ID common.IDStr

type Content struct {
	Title    string      `json:",omitempty" bson:",omitempty"`
	Summary  string      `json:",omitempty" bson:",omitempty"`
	TypeKey  TypeKey     `json:",omitempty" bson:",omitempty"`
	Data     string      `json:",omitempty" bson:",omitempty"`
	Embedded []Content   `json:",omitempty" bson:",omitempty"` // in particular: URLs, images, etc.
	Tags     []tags.Item `json:",omitempty" bson:",omitempty"`
}

type Item struct {
	ID        ID          `json:",omitempty" bson:"_id,omitempty"`
	IssuedID  ns.ID       `json:",omitempty" bson:",omitempty"`
	Content   Content     `json:",inline"    bson:",inline"`
	OwnerID   auth.ID     `json:",omitempty" bson:",omitempty"`
	ViewerID  auth.ID     `json:",omitempty" bson:",omitempty"`
	History   vcs.History `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time   `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time  `json:",omitempty" bson:",omitempty"`
}

type Operator interface {
	Save(Item, *crud.Options) (*Item, error)
	Remove(ID, *crud.Options) error
	Read(ID, *crud.Options) (*Item, error)
	List(*crud.Options) ([]Item, error) // in particular: selected owned, tagged, untagged, containing string, etc.
	Tags(*crud.Options) (tags.StatMap, error)
	// StatMap(*crud.Options) (common.Map, error) // in particular: selected, grouped, etc.

	HasTag(tag tags.Item) (selectors.Term, error)
	AddParent(tags []tags.Item, id ID) ([]tags.Item, error)
	HasParent(id ID) (selectors.Term, error)
}
