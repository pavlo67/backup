package catalogue_files

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/selectors"
)

var _ db.Cleaner = &catalogueFiles{}

const onClean = "on filesFS.Clean()"

func (filesOp *catalogueFiles) Clean(term *selectors.Term) error {
	//if err := os.RemoveAll(filesOp.basePath); err != nil {
	//	return errors.Wrapf(err, onClean+": removing %s", filesOp.basePath)
	//}

	return common.ErrNotImplemented
}
