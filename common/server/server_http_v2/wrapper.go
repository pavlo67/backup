package server_http_v2

import (
	"github.com/julienschmidt/httprouter"
)

type HandlerHTTP = httprouter.Handle
type WrapperHTTP func(op OperatorV2, serverPath string, data interface{}) (method, path string, h HandlerHTTP, err error)
