package notebook_html

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cbroglie/mustache"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
	"github.com/pavlo67/tools/components/views/views_html"
)

var _ Operator = &notebookHTML{}

// should be thread-safe
type notebookHTML struct {
	htmlTemplate string
	pagesConfig  server_http.Config
	restConfig   server_http.Config
}

const onNew = "on notebookHTML.New(): "

func New(htmlTemplate string, pagesConfig, restConfig server_http.Config) (Operator, error) {
	if strings.TrimSpace(htmlTemplate) == "" {
		return nil, errors.New("no htmlTemplate to render pages")
	}

	method, urlStr, err := pagesConfig.EP(notebook.IntefaceKeyHTMLView, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = pagesConfig.EP(notebook.IntefaceKeyHTMLTagged, []string{"id"}, false)

	// l.Infof("%s %s", method, urlStr)

	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLList, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLList, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = pagesConfig.EP(notebook.IntefaceKeyHTMLEdit, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLEdit, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLEdit, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	return &notebookHTML{
		htmlTemplate: htmlTemplate,
		pagesConfig:  pagesConfig,
		restConfig:   restConfig,
	}, nil
}

// TODO!!! look at https://github.com/kataras/blocks

func (htmlOp *notebookHTML) HTMLMessage(errs errors.Error) string {
	return errs.Error()
}

func (htmlOp *notebookHTML) HTMLPage(title, htmlHeader, htmlIndex, htmlContent, htmlMessage string) server.Response {
	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(fmt.Sprintf("%s %s %s %s %s", title, htmlHeader, htmlMessage, htmlIndex, htmlContent)),
		MIMEType: "text/html; charset=utf-8",
	}
}

func (htmlOp *notebookHTML) HTMLError(httpStatus int, htmlError string) (server.Response, error) {
	context := map[string]string{
		"title":   "нотатник: помилка",
		"header":  "Помилка",
		"content": htmlError,
	}

	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   httpStatus,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err
}

func (htmlOp *notebookHTML) HTMLRoot(htmlHello, htmlTags string) (server.Response, error) {

	context := map[string]string{
		"title":   "нотатник: вхід",
		"header":  "Вхід",
		"content": htmlHello + "\n\n<p>" + htmlTags,
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err
}

const onHTMLView = "on notebookHTML.HTMLView(): "

func (htmlOp *notebookHTML) HTMLView(r *records.Item, children []records.Item, message string) (server.Response, error) {
	context := map[string]string{
		"title":   "нотатник: " + r.Content.Title,
		"header":  r.Content.Title,
		"message": message,
		"content": views_html.HTMLViewTable(dataFields, DataFromRecord(r), nil),
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err
}

const onHTMLEdit = "on notebookHTML.HTMLEdit(): "

func (htmlOp *notebookHTML) HTMLEdit(r *records.Item, children []records.Item, message string) (server.Response, error) {
	formID := "nb_edit_" + strconv.FormatInt(time.Now().Unix(), 10) + "_"

	updateFields := append(
		dataFields,
		views_html.Field{
			"update",
			"зберегти зміни",
			"submit",
			nil,
			map[string]string{"class": "ut"},
		},
		//views_html.Field{
		//	"update",
		//	"зберегти зміни",
		//	"button",
		//	nil,
		//	map[string]string{"class": "ut", "onclick": `getData("` + formID + `")`},
		//},
	)

	context := map[string]string{
		"title":   "редаґування нотатки: " + r.Content.Title,
		"header":  "редаґування нотатки: " + r.Content.Title,
		"message": message,
		"content": views_html.HTMLEditTable(updateFields, formID, "/save", DataFromRecord(r), nil),
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err

}

func (htmlOp *notebookHTML) HTMLList(tag tags.Item, tagged []records.Item) (string, error) {
	return fmt.Sprintf("%s / %#v", tag, tagged), nil

}

const onHTMLTags = "on tagsHTML.HTMLTags(): "

func (htmlOp *notebookHTML) HTMLTags(tss tags.Stats) (string, error) {
	htmlTags := "\n\n<ul>\n"

	for i, ts := range tss {
		tag := strings.TrimSpace(string(ts.Tag))
		if tag == "" {
			l.Errorf("empty tag %d in %#v", i, tss)
			continue
		}

		method, urlStr, err := htmlOp.pagesConfig.EP(notebook.IntefaceKeyHTMLTagged, []string{tag}, false)

		l.Infof("%s %s", method, urlStr)

		if err != nil || urlStr == "" {
			l.Errorf("can't server_http.EP(%#v, notebook.IntefaceKeyHTMLTags,nil,false), got %s, %s, %s", htmlOp.pagesConfig, method, urlStr, err)
			continue
		}

		htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	}

	htmlTags += "</ul>\n\n"

	return htmlTags, nil
}
