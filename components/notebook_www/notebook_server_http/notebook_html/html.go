package notebook_html

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/components/notebook_www"

	"github.com/pavlo67/data_exchange/components/tags"
	"github.com/pavlo67/tools/common/views/views_html"
	"github.com/pavlo67/tools/entities/records"
)

var _ Operator = &notebookHTML{}

type notebookHTML struct {
	epCreate string
	epView   server_http.Get1
	epTagged server_http.Get1
}

const onNew = "on notebookHTML.New(): "

func New(pagesConfig server_http.ConfigPages) (Operator, error) { // , restConfig

	epCreate, err := server_http.CheckGet0(pagesConfig, notebook_www.IntefaceKeyHTMLCreate, false)
	if err != nil {
		return nil, err
	}

	epView, err := server_http.CheckGet1(pagesConfig, notebook_www.IntefaceKeyHTMLView, false)
	if err != nil {
		return nil, err
	}

	epTagged, err := server_http.CheckGet1(pagesConfig, notebook_www.IntefaceKeyHTMLTagged, false)
	if err != nil {
		return nil, err
	}

	return &notebookHTML{
		epCreate: epCreate,
		epView:   epView,
		epTagged: epTagged,
	}, nil
}

// TODO!!! look at https://github.com/kataras/blocks

func (htmlOp *notebookHTML) FragmentsView(r *records.Item, children []records.Item, message string, identity *auth.Identity) (server_http.Fragments, error) {
	fragments := server_http.Fragments{
		"title":   r.Content.Title,
		"header":  r.Content.Title,
		"message": message,
		"content": views_html.HTMLViewTable(dataFields, DataFromRecord(r), nil),
	}

	return fragments, nil
}

const onHTMLEdit = "on notebookHTML.FragmentsEdit(): "

func (htmlOp *notebookHTML) FragmentsEdit(r *records.Item, children []records.Item, message string, identity *auth.Identity) (server_http.Fragments, error) {
	formID := "nb_edit_" + strconv.FormatInt(time.Now().Unix(), 10) + "_"

	var title, header, action string
	var dataFromRecord server_http.Fragments
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
			server_http.Fragments{"class": "ut"},
		},
	)

	fragments := server_http.Fragments{
		"title":   title,
		"header":  header + title,
		"message": message,
		"content": views_html.HTMLEditTable(updateFields, formID, "/save", dataFromRecord, nil),
	}

	return fragments, nil
}

func (htmlOp *notebookHTML) FragmentsListTagged(tag tags.Item, tagged []records.Item, identity *auth.Identity) (server_http.Fragments, error) {
	htmlList := htmlOp.HTMLFiles(tagged, identity)

	fragments := server_http.Fragments{
		"title":   "все з теґом '" + tag + "'",
		"header":  "Все з теґом '" + tag + "'",
		"content": htmlList,
	}
	return fragments, nil
}

func (htmlOp *notebookHTML) HTMLTags(tagsStatMap tags.StatMap, identity *auth.Identity) string {

	tss := tagsStatMap.List(true)

	htmlTags := "\n\n<ul>\n"

	for i, ts := range tss {
		tag := strings.TrimSpace(string(ts.Tag))
		if tag == "" {
			l.Errorf("empty tag %d in %#v", i, tss)
			continue
		}

		urlStr, err := htmlOp.epTagged(tag)
		if err != nil || urlStr == "" {
			l.Errorf("can't htmlOp.epTagged(%s), got %s, %s", tag, urlStr, err)
			continue
		}

		htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	}
	htmlTags += "</ul>\n\n"

	return htmlTags

}

func (htmlOp *notebookHTML) HTMLIndex(identity *auth.Identity) string {
	htmlIndex := `<div style="padding:5px;margin: 15px 0 10px 10px;width:200px;float:right;">`

	urlStr := htmlOp.epCreate
	htmlIndex += fmt.Sprintf(`<li>[<a href="%s">новий запис</a>]</li>`, urlStr)
	htmlIndex += "</div>\n\n"

	return htmlIndex

}

func (htmlOp *notebookHTML) HTMLFiles(recordItems []records.Item, identity *auth.Identity) string {
	if len(recordItems) < 1 {
		return "нема записів"
	}

	var htmlRecords string

	for _, r := range recordItems {
		details := `<table class="border" style="padding:3px;margin: 0 0 10px 10px;width:150px;" align=right>` +
			"<tr><td>" + HTMLAuthor(&r, identity) + "</td></tr>\n" +
			"<tr><td>" + HTMLTags(r.Tags, r.ViewerNSS, r.OwnerNSS, htmlOp.epTagged, "<br>- ") + "</td></tr>\n" +
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
			listDelimiter + "\n"
	}

	return htmlRecords

}

func HTMLHidden(html string) string {
	return `<div style="position:absolute;visibility:hidden;">` + html + "</div>\n"
}
