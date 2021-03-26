package server_http

type WrapperHTTPKey string

const WrapperHTTPREST WrapperHTTPKey = "rest"
const WrapperHTTPPage WrapperHTTPKey = "page"

// const WrapperHTTPFiles WrapperHTTPKey = "files"

type OperatorV2 interface {
	Handle(key EndpointKey, serverPath string, wrapperHTTPKey WrapperHTTPKey, data interface{}) error
	HandleFiles(key EndpointKey, serverPath string, staticPath StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
