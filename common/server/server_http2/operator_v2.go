package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"
)

const OnRequestMiddlewareInterfaceKey joiner.InterfaceKey = "server_http_on_request_middleware"
const InterfaceKey joiner.InterfaceKey = "server_http"

type PathParams map[string]string

type WrapperHTTPKey string

const WrapperHTTPREST WrapperHTTPKey = "rest"
const WrapperHTTPPage WrapperHTTPKey = "page"
const WrapperHTTPFiles WrapperHTTPKey = "files"

type OnRequestMiddleware interface {
	Identity(r *http.Request) (*auth.Identity, error)
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type OperatorV2 interface {
	Handle(key EndpointKey, serverPath string, wrapperHTTPKey WrapperHTTPKey, data interface{}) error
	HandleFiles(key EndpointKey, serverPath string, staticPath StaticPath) error
	Start() error
	Addr() (port int, https bool)
}
