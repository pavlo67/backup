package notebook_server_http

import (
	"github.com/pavlo67/common/common/crud"

	"github.com/pavlo67/tools/components/records"
)

func prepareRecord(id records.ID, options *crud.Options) (*records.Item, []records.Item, error) {
	r, err := recordsOp.Read(id, options)
	if err != nil {
		return r, nil, err
	}

	selector, err := recordsOp.HasParent(id)
	if err != nil {
		return r, nil, err
	}

	options = options.WithSelector(selector)
	children, err := recordsOp.List(options)
	return r, children, err
}
