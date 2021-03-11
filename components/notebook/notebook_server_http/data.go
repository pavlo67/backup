package notebook_server_http

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/tools/components/records"
)

func prepareRecord(id records.ID, options *crud.Options) (*records.Item, []records.Item, error) {
	r, err := recordsOp.Read(id, options)
	if err != nil {
		return r, nil, err
	}

	selectorParent := selectors.Term{
		Key:    records.HasParent,
		Values: []string{string(id)},
	}

	options = options.WithSelector(selectorParent)
	children, err := recordsOp.List(options)
	return r, children, err
}
