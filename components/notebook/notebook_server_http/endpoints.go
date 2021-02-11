package notebook_server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
)

var Endpoints = server_http.Endpoints{
	notebook.IntefaceKeyRoot: rootEndpoint,
}

var rootEndpoint = server_http.Endpoint{
	Method: "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {
		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte("мама мила раму!!!"),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}
