package views_html

import "strconv"

var num int

func HTMLOpenClose(title, content string, visible bool) string {
	num++

	id := strconv.Itoa(num)
	imageID := "plus_minus_img" + id
	contentID := "plus_minus_content" + id

	var htmlContent string

	if visible {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + ImgMinus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `">` + content + `</div>`
	} else {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + ImgPlus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `" style="visibility:hidden;position:absolute;">` + content + `</div>`

	}

	return htmlContent
}
