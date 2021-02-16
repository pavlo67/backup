package notebook_server_http

import (
	"fmt"
	"html"
	"net/http"

	"github.com/pavlo67/common/common/server"

	"github.com/pavlo67/tools/components/views/views_html"
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

func HTMLEditTable(fields []views_html.Field, formID string, data map[string]string, values map[string]views_html.SelectString) string { // ,
	// frontOps map[string]views_html.Operator, rView auth.ID, publicChanges bool
	//if data == nil {
	//	data = map[string]string{}
	//}
	//if values == nil {
	//	values = map[string]views_html.SelectString{}
	//}
	//
	//values["visibility"], data["visibility"] = dataView(user, rView, publicChanges)

	var editHTML, titleHTML, resHTML string
	for _, f := range fields {
		titleHTML, resHTML = views_html.FieldEdit("edit_note_"+formID, f, data, values) // , frontOps

		//if resHTML == "" && f.Params[views_html.NotEmptyKey] == true {
		//	continue
		//}

		if titleHTML != "" {
			titleHTML = "<small>" + titleHTML + ":</small> \n"
		}
		editHTML += `<tr id="div_` + html.EscapeString(formID+f.Key) + `"><td>` + "\n" + titleHTML + "</td><td>" + resHTML + "</td></tr>"

	}

	return `<table width="100%">` + editHTML + "</table>\n"
}
