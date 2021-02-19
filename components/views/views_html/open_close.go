package views_html

import (
	"strconv"

	"github.com/pavlo67/tools/components/notebook/notebook_html"
)

var num int

func HTMLOpenClose(title, content string, visible bool) string {
	num++

	id := strconv.Itoa(num)
	imageID := "plus_minus_img" + id
	contentID := "plus_minus_content" + id

	var htmlContent string

	if visible {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + notebook_html.ImgMinus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `">` + content + `</div>`
	} else {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + notebook_html.ImgPlus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `" style="visibility:hidden;position:absolute;">` + content + `</div>`

	}

	return htmlContent
}
