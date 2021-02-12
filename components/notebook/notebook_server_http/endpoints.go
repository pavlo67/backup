package notebook_server_http

import (
	"net/http"

	"github.com/pavlo67/tools/components/formatter"

	"github.com/pavlo67/tools/components/records"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
)

var Endpoints = server_http.Endpoints{
	rootEndpoint,
	editEndpoint,
	viewEndpoint,
}

var rootEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyRoot,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {
		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte("мама мила раму!!!"),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var viewEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			l.Error(err)
		}

		htmlStr, err := formatterRecordsOp.Format(r, formatter.Full)
		if err != nil {
			l.Error(err)
		}

		return ResponseHTMLOk(0, htmlStr), nil
	},
}

var editEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			l.Error(err)
		}

		htmlStr, err := formatterRecordsOp.Format(r, formatter.Edit)
		if err != nil {
			l.Error(err)
		}

		return ResponseHTMLOk(0, htmlStr), nil
	},
}
