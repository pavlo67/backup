package server_http_v2

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
)

type EndpointsRESTSettled map[server_http.EndpointKey]EndpointRESTSettled

type EndpointRESTSettled struct {
	Path     string
	Tags     []string `json:",omitempty"`
	Produces []string `json:",omitempty"`
	EndpointREST
}

type WorkerHTTPv2 func(OperatorV2, *http.Request, server_http.PathParams, *auth.Identity) (server.Response, error)

type EndpointREST struct {
	server_http.EndpointDescription
	WorkerHTTPv2
}
