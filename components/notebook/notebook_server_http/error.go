package notebook_server_http

import (
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http/notebook_html"
)

func errorPage(httpStatus int, notebookHTMLOp notebook_html.Operator, err error, publicDetails string, req *http.Request) (server.Response, error) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	htmlPage, errRender := notebookHTMLOp.CommonPage(
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

	if req != nil {
		err = errors.CommonError(fmt.Errorf(" on [%s %s]", req.Method, req.URL), publicDetails, err)
	} else {
		err = errors.CommonError(publicDetails, err)
	}

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err
}
