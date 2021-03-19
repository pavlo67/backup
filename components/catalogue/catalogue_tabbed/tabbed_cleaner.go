package catalogue_tabbed

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/selectors"
)

var _ db.Cleaner = &catalogueTabbed{}

const onClean = "on filesFS.Clean()"

func (catalogueOp *catalogueTabbed) Clean(term *selectors.Term) error {
	return common.ErrNotImplemented
}
