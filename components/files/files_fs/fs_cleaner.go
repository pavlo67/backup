package files_fs

import (
	"os"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
)

var _ crud.Cleaner = &filesFS{}

const onClean = "on filesFS.Clean()"

func (filesOp *filesFS) Clean(opts *crud.Options) error {
	for bucketID, basePath := range filesOp.buckets {
		if err := os.RemoveAll(basePath); err != nil {
			return errors.Wrapf(err, onClean+": removing %s --> %s", bucketID, basePath)
		}
	}

	return nil
}
