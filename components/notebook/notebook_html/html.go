package notebook_html

import (
	"encoding/json"
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

	// TODO!!! check by single call with required endpoints list

	method, urlStr, err := pagesConfig.EP(notebook.IntefaceKeyHTMLView, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = pagesConfig.EP(notebook.IntefaceKeyHTMLTagged, []string{"id"}, false)
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

func pageError(err, errRender error, req *http.Request) error {
	var errs []interface{}

	if err != nil {
		errs = []interface{}{err}
	}
	if errRender != nil {
		errs = append(errs, errRender)
	}

	if len(errs) < 1 {
		return nil
	}

	if req != nil {
		return errors.CommonError(append([]interface{}{fmt.Errorf("on %s %s", req.Method, req.URL)}, errs...)...)
	}

	return errors.CommonError(errs...)
}

func (htmlOp *notebookHTML) HTMLError(httpStatus int, err error, publicDetails string, req *http.Request) (server.Response, error) {
	context := map[string]string{
		"title":   "нотатник: помилка",
		"header":  "Помилка",
		"content": "На жаль, виникла помилка (" + publicDetails + ")",
	}

	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	htmlPage, errRender := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   httpStatus,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, pageError(err, errRender, req)
}

func (htmlOp *notebookHTML) HTMLRoot(htmlHello string, tagsStatMap tags.StatMap) (server.Response, error) {

	context := map[string]string{
		"title":   "нотатник: вхід",
		"header":  "Вхід",
		"content": htmlHello + "\n\n<p>" + htmlOp.htmlTags(tagsStatMap),
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, pageError(err, errRender, req)
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

func (htmlOp *notebookHTML) HTMLTagged(tag tags.Item, tagged []records.Item) (server.Response, error) {
	htmlList, err := json.Marshal(tagged)

	context := map[string]string{
		"title":   "нотатник: все з теґом '" + tag + "'",
		"header":  "Все з теґом '" + tag + "'",
		"content": string(htmlList),
	}

	htmlPage, err := mustache.Render(htmlOp.htmlTemplate, context)

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err

}

const onHTMLTags = "on tagsHTML.HTMLTags(): "

func (htmlOp *notebookHTML) HTMLTags(tagsStatMap tags.StatMap) (server.Response, error) {

	title := "нотатник: теґи"
	htmlHeader := "Теґи"

	return notebookHTMLOp.HTMLPage(title, htmlHeader, "", htmlTags, errs.Error()), nil

}

func (htmlOp *notebookHTML) htmlTags(tagsStatMap tags.StatMap) string {

	tss := tagsStatMap.List(true)

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
			l.Errorf("can't server_http.EP(%#v, notebook.InterfaceKeyHTMLTags,nil,false), got %s, %s, %s", htmlOp.pagesConfig, method, urlStr, err)
			continue
		}

		htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	}
	htmlTags += "</ul>\n\n"

	return htmlTags

}
