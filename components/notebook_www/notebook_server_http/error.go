package notebook_server_http

import (
	"net/http"

	server_http "github.com/pavlo67/tools/common/server/server_http2"

	"github.com/pavlo67/tools/components/notebook_www/notebook_server_http/notebook_html"
)

func errorPage(httpStatus int, notebookHTMLOp notebook_html.Operator, err error, publicDetails string, req *http.Request) (server_http.ResponsePage, error) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	context, errRender := notebookHTMLOp.CommonPage(
		"помилка",
		"",
		"",
		publicDetails,
		"",
		"",
	)
	if errRender != nil {
		l.Error(err)
	}

	//
	//if req != nil {
	//	err = errors.CommonError(fmt.Errorf(" on [%s %s]", req.Method, req.URL), publicDetails, err)
	//} else {
	//	err = errors.CommonError(publicDetails, err)
	//}
	//
	return server_http.ResponsePage{
		Status:    http.StatusOK,
		Fragments: context,
	}, err
}
