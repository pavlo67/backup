package catalogue_tabbed

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
)

var _ crud.Cleaner = &catalogueTabbed{}

const onClean = "on filesFS.Clean()"

func (catalogueOp *catalogueTabbed) Clean(opts *crud.Options) error {
	return common.ErrNotImplemented
}
