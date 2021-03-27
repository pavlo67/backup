package wrapper_page

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pavlo67/common/common/logger"

	"github.com/cbroglie/mustache"
	"github.com/julienschmidt/httprouter"
	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

// Page ----------------------------------------------------------------------------------------------------

type CommonFragments interface {
	Save(key, fragment string) error
	Set(fragments server_http.Fragments) (server_http.Fragments, error)
}

func WrapperHTTPPage(htmlTemplate string, commonFragments CommonFragments, l logger.Operator) server_http.WrapperHTTP {
	return func(serverOpV2 server_http.OperatorV2, serverPath string, data interface{}) (string, string, server_http.HandlerHTTP, error) {
		var ep *server_http.EndpointPage

		switch v := data.(type) {
		case server_http.EndpointPage:
			ep = &v
		case *server_http.EndpointPage:
			ep = v
		}

		if ep == nil {
			return "", "", nil, fmt.Errorf("wrong data for WrapperHTTPPage: %#v", data)
		}

		handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
			//options, err := s.onRequest.Identity(r)
			//if err != nil {
			//	l.Error(err)
			//}
			var identity *auth.Identity

			var params server_http.PathParams
			if len(paramsHR) > 0 {
				params = server_http.PathParams{}
				for _, p := range paramsHR {
					params[p.Key] = p.Value
				}
			}

			responseData, err := ep.WorkerHTTPPage(serverOpV2, r, params, identity)
			if err != nil {
				l.Error("on ep.WorkerHTTPPage(): ", err)
			}

			if responseData.Status > 0 {
				w.WriteHeader(responseData.Status)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			var fragments server_http.Fragments
			if commonFragments == nil {
				fragments = responseData.Fragments
			} else if fragments, err = commonFragments.Set(fragments); err != nil {
				l.Error("on commonFragments(): ", err)
			}

			responseBody, err := mustache.Render(htmlTemplate, fragments)
			if err != nil {
				// TODO!!!
				l.Error("on mustache.Render(): ", err)
			}

			if _, err := w.Write([]byte(responseBody)); err != nil {
				l.Error("can't write response", err)
			}
		}

		method := strings.ToUpper(ep.Method)
		path := ep.PathTemplate(serverPath)

		return method, path, handler, nil
	}
}
