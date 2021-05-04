package server_http_v2

import "github.com/pavlo67/common/common/server/server_http"

type WrapperHTTPKey string

const WrapperHTTPREST WrapperHTTPKey = "rest"
const WrapperHTTPPage WrapperHTTPKey = "page"

// const WrapperHTTPFiles WrapperHTTPKey = "files"

type OperatorV2 interface {
	Handle(key server_http.EndpointKey, serverPath string, wrapperHTTPKey WrapperHTTPKey, data interface{}) error
	HandleFiles(key server_http.EndpointKey, serverPath string, staticPath server_http.StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
