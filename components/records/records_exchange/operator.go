package records_exchange

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
)

type ID common.IDStr

type Description struct {
}

type Record struct {
}

type Package struct {
	ID          ID          `json:"id"          bson:"_id"`
	Description Description `json:"description" bson:"description"`
	Records     []Record    `json:"records"     bson:"records"`
}

type Operator interface {
	Exchange(Package, *crud.Options) (Package, error)
}
