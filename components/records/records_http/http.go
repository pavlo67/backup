package records_http

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

var _ records.Operator = &recordsHTTP{}

type recordsHTTP struct {
	pagesConfig server_http.Config
	restConfig  server_http.Config
}

const onNew = "on recordsHTTP.New()"

func New(pagesConfig, restConfig server_http.Config) (records.Operator, error) {
	// TODO!!! check endpoints in config

	recordsOp := recordsHTTP{
		pagesConfig: pagesConfig,
		restConfig:  restConfig,
	}

	return &recordsOp, nil
}

func (recordsOp *recordsHTTP) Save(records.Item, *crud.Options) (*records.Item, error) {
	//ep := recordsOp.pagesConfig.EndpointsSettled[records.IntefaceKeySetCreds]
	//serverURL := recordsOp.pagesConfig.Host + recordsOp.pagesConfig.Port + ep.Path
	//
	//requestBody, err := json.Marshal(toSet)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onrecordsenticate+": can't marshal toSet(%#v)", toSet)
	//}
	//
	//var creds *records.Creds
	//if err := server_http.Request(serverURL, ep, requestBody, creds, &crud.Options{Identity: &records.Identity{ID: recordsID}}, l); err != nil {
	//	return nil, err
	//}
	//
	//return creds, nil

	return nil, common.ErrNotImplemented
}

func (recordsOp *recordsHTTP) Remove(records.ID, *crud.Options) error {
	return common.ErrNotImplemented
}

func (recordsOp *recordsHTTP) Read(records.ID, *crud.Options) (*records.Item, error) {
	return nil, common.ErrNotImplemented
}

func (recordsOp *recordsHTTP) List(*crud.Options) ([]records.Item, error) {
	return nil, common.ErrNotImplemented
}

func (recordsOp *recordsHTTP) Tags(*crud.Options) (tags.StatMap, error) {
	return nil, common.ErrNotImplemented
}

func (recordsOp *recordsHTTP) AddParent(tags []tags.Item, id records.ID) ([]tags.Item, error) {
	return nil, common.ErrNotImplemented
}
