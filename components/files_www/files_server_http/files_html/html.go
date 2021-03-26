package files_html

import (
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/components/files_www"

	"github.com/pavlo67/tools/common/views/views_html"
	"github.com/pavlo67/tools/entities/records"
)

var _ Operator = &filesHTML{}

type filesHTML struct {
	epCreate string
	epView   server_http.Get1
}

const onNew = "on filesHTML.New(): "

func New(pagesConfig server_http.ConfigPages) (Operator, error) { // , restConfig

	epCreate, err := server_http.CheckGet0(pagesConfig, files_www.IntefaceKeyHTMLCreate, false)
	if err != nil {
		return nil, err
	}

	epView, err := server_http.CheckGet1(pagesConfig, files_www.IntefaceKeyHTMLView, false)
	if err != nil {
		return nil, err
	}

	return &filesHTML{
		epCreate: epCreate,
		epView:   epView,
	}, nil
}

// TODO!!! look at https://github.com/kataras/blocks

func (htmlOp *filesHTML) CommonPage(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) (map[string]string, error) {

	if htmlError = strings.TrimSpace(htmlError); htmlError != "" {
		htmlError = "На жаль, виникла помилка:-(\n<p>" + htmlError
	}

	context := map[string]string{
		"title":   title,
		"header":  htmlHeader,
		"message": htmlMessage,
		"error":   htmlError,
		"index":   htmlIndex,
		"content": htmlContent,
	}

	return context, nil
}

func (htmlOp *filesHTML) View(r *records.Item, children []records.Item, message string, identity *auth.Identity) (map[string]string, error) {
	context := map[string]string{
		"title":   r.Content.Title,
		"header":  r.Content.Title,
		"message": message,
		"content": views_html.HTMLViewTable(dataFields, DataFromRecord(r), nil),
	}

	return context, nil
}

const onHTMLEdit = "on filesHTML.Edit(): "

func (htmlOp *filesHTML) Edit(r *records.Item, children []records.Item, message string, identity *auth.Identity) (map[string]string, error) {
	formID := "nb_edit_" + strconv.FormatInt(time.Now().Unix(), 10) + "_"

	var title, header, action string
	var dataFromRecord map[string]string
	if r == nil {
		header = "Створення запису"
		action = "зберегти запис"
	} else {
		header = "Редаґування: "
		title = r.Content.Title
		dataFromRecord = DataFromRecord(r)
		action = "зберегти зміни"
	}

	updateFields := append(
		dataFields,
		views_html.Field{
			"update",
			action,
			"submit",
			nil,
			map[string]string{"class": "ut"},
		},
	)

	context := map[string]string{
		"title":   title,
		"header":  header + title,
		"message": message,
		"content": views_html.HTMLEditTable(updateFields, formID, "/save", dataFromRecord, nil),
	}

	return context, nil
}

func (htmlOp *filesHTML) HTMLFiles(recordItems []records.Item, identity *auth.Identity) string {
	if len(recordItems) < 1 {
		return "нема записів"
	}

	var htmlRecords string

	for _, r := range recordItems {
		details := `<table class="border" style="padding:3px;margin: 0 0 10px 10px;width:150px;" align=right>` +
			`</table>` +
			"<p>" + r.Content.Summary
		// + HTMLFiles(r.Links, pxPreview)

		name := strings.TrimSpace(r.Content.Title)
		if name == "" {
			name = "..."
		}

		urlStr, err := htmlOp.epView(string(r.ID))
		if err != nil || urlStr == "" {
			l.Errorf("can't htmlOp.epView(%s), got %s, %s", r.ID, urlStr, err)
		}

		htmlRecords += `<li><a href="` + urlStr + `">` + name + "</a></li>\n" +
			"<br>" + details + // HTMLHidden(details) +
			"\n"
	}

	return htmlRecords

}
