package environments

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/joiner"
)

type ID common.IDStr

type Operator interface {
	Get(ID) (joiner.Operator, error)
}
