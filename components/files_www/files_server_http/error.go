package files_server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/errors"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/tools/components/files_www/files_server_http/files_html"
)

func errorPage(httpStatus int, filesHTMLOp files_html.Operator, err error, publicDetails string, req *http.Request) (server_http.ResponsePage, error) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	context, errRender := filesHTMLOp.CommonPage(
		"помилка",
		"",
		"",
		publicDetails,
		"",
		"",
	)
	if errRender != nil {
		l.Error(publicDetails, " ", err)
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
	}, errors.CommonError(publicDetails, " ", err)
}
