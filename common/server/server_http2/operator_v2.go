package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/common/common/server"
)

const OnRequestMiddlewareInterfaceKey joiner.InterfaceKey = "server_http_on_request_middleware"
const InterfaceKey joiner.InterfaceKey = "server_http"

type PathParams map[string]string
type WorkerHTTP func(OperatorV2, *http.Request, PathParams, *auth.Identity) (server.Response, error)
type PreparatorHTTP func(interface{}) (WorkerHTTP, error)

type OnRequestMiddleware interface {
	Identity(r *http.Request) (*auth.Identity, error)
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type OperatorV2 interface {
	HandleEndpoint(key EndpointKey, serverPath string, endpoint Endpoint) error
	HandleFiles(key EndpointKey, serverPath string, staticPath StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
