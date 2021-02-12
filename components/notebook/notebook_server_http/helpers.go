package notebook_server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/server"
)

func ResponseHTMLOk(status int, htmlData string) server.Response {
	if status <= 0 {
		status = http.StatusOK
	}

	return server.Response{
		Status:   status,
		Data:     []byte(htmlData),
		MIMEType: "text/html; charset=utf-8",
	}
}
