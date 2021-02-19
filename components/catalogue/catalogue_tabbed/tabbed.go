package catalogue_tabbed

import (
	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/tools/components/catalogue"
	"github.com/pavlo67/tools/components/files"
)

var _ catalogue.Operator = &catalogueTabbed{}

type catalogueTabbed struct {
	filesOp files.Operator
}

const onNew = "on catalogueTabbed.New(): "

func New(filesOp files.Operator) (catalogue.Operator, crud.Cleaner, error) {
	if filesOp == nil {
		return nil, nil, errors.New(onNew + ": no files.Operator")
	}

	catalogueOp := catalogueTabbed{
		filesOp: filesOp,
	}

	return &catalogueOp, &catalogueOp, nil
}

const onSave = "on catalogueTabbed.Save()"

func (catalogueOp *catalogueTabbed) Save(bucketID files.BucketID, item catalogue.Item, options *crud.Options) (string, error) {
	return "", common.ErrNotImplemented
}

const onRead = "on catalogueTabbed.Read()"

func (catalogueOp *catalogueTabbed) Read(bucketID files.BucketID, path string, options *crud.Options) (*catalogue.Item, error) {
	return nil, common.ErrNotImplemented
}

const onRemove = "on catalogueTabbed.Remove()"

func (catalogueOp *catalogueTabbed) Remove(bucketID files.BucketID, path string, options *crud.Options) error {
	return common.ErrNotImplemented
}

const onList = "on catalogueTabbed.Items()"

func (catalogueOp *catalogueTabbed) List(bucketID files.BucketID, path string, depth int, options *crud.Options) ([]catalogue.Item, error) {
	return nil, common.ErrNotImplemented
}

const onStat = "on catalogueTabbed.Stat()"

func (catalogueOp *catalogueTabbed) Stat(bucketID files.BucketID, path string, depth int, options *crud.Options) (*catalogue.Item, error) {
	return nil, common.ErrNotImplemented
}
