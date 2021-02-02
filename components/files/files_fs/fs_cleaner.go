package files_fs

import (
	"github.com/pavlo67/common/common/crud"
)

var _ crud.Cleaner = &filesFS{}

const onClean = "on filesFS.Clean()"

func (filesOp *filesFS) Clean(opts *crud.Options) error {
	return nil
}
