package notebook_server_http

import (
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common/server"
)

const pageTemplate = `<html>
<head>
<title>%s</title>
</head>
<body>
<div><b>%s</b></div>
<div>%s</div>
<div align="right">%s</div>
<div>%s</div>
</body>
</html>
`

func HTMLPage(title, htmlHeader, htmlIndex, htmlContent, htmlMessage string) server.Response {
	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(fmt.Sprintf(pageTemplate, title, htmlHeader, htmlMessage, htmlIndex, htmlContent)),
		MIMEType: "text/html; charset=utf-8",
	}
}
