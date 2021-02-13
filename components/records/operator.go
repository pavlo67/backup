package records

import (
	"sort"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/tools/components/ns"
	"github.com/pavlo67/tools/components/vcs"
)

type ID common.IDStr

type Content struct {
	Title    string    `json:",omitempty" bson:",omitempty"`
	Summary  string    `json:",omitempty" bson:",omitempty"`
	TypeKey  TypeKey   `json:",omitempty" bson:",omitempty"`
	Data     string    `json:",omitempty" bson:",omitempty"`
	Embedded []Content `json:",omitempty" bson:",omitempty"` // in particular: URLs, images, etc.
	Tags     []string  `json:",omitempty" bson:",omitempty"`
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
	Tags(*crud.Options) (TagsStat, error)
	// Stat(*crud.Options) (common.Map, error) // in particular: selected, grouped, etc.

	HasTag(tag string) (selectors.Term, error)
	AddParent(tags []string, id ID) ([]string, error)
	HasParent(id ID) (selectors.Term, error)
}

type TagsStat map[string]int64
type TagStatItem struct {
	Tag   string
	Count int64
}
type TagsStatList []TagStatItem

func (ts TagsStat) StatList(sortByCount bool) TagsStatList {
	var tagsStatList TagsStatList
	for t, c := range ts {
		tagsStatList = append(tagsStatList, TagStatItem{t, c})
	}

	if sortByCount {
		sort.Slice(tagsStatList, func(i, j int) bool { return tagsStatList[i].Count >= tagsStatList[j].Count })
	} else {
		sort.Slice(tagsStatList, func(i, j int) bool { return tagsStatList[i].Tag <= tagsStatList[j].Tag })
	}

	return tagsStatList
}
