package notebook_server_http

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/tools/components/records"
)

func prepareRecord(id records.ID, identity *auth.Identity) (*records.Item, []records.Item, error) {
	r, err := recordsOp.Read(id, identity)
	if err != nil {
		return r, nil, err
	}

	selectorParent := selectors.Term{
		Key:    records.HasParent,
		Values: []string{string(id)},
	}

	children, err := recordsOp.List(&selectorParent, identity)
	return r, children, err
}
