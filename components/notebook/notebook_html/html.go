package notebook_html

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cbroglie/mustache"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
	"github.com/pavlo67/tools/components/views/views_html"
)

var _ Operator = &notebookHTML{}

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

func (htmlOp *notebookHTML) CommonPage(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) (string, error) {

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

	return mustache.Render(htmlOp.htmlTemplate, context)
}

func (htmlOp *notebookHTML) View(r *records.Item, children []records.Item, message string, options *crud.Options) (string, error) {
	context := map[string]string{
		"title":   r.Content.Title,
		"header":  r.Content.Title,
		"message": message,
		"content": views_html.HTMLViewTable(dataFields, DataFromRecord(r), nil),
	}

	return mustache.Render(htmlOp.htmlTemplate, context)
}

const onHTMLEdit = "on notebookHTML.Edit(): "

func (htmlOp *notebookHTML) Edit(r *records.Item, children []records.Item, message string, options *crud.Options) (string, error) {
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

	return mustache.Render(htmlOp.htmlTemplate, context)
}

func (htmlOp *notebookHTML) ListTagged(tag tags.Item, tagged []records.Item, options *crud.Options) (string, error) {
	htmlList, errRenderRecords := htmlOp.HTMLRecords(tagged, options)
	if errRenderRecords != nil {
		errorID := strconv.FormatInt(time.Now().UnixNano(), 10)
		l.Error(errorID, " / ", errRenderRecords)
		return htmlOp.CommonPage(
			"помилка",
			"",
			"",
			"при htmlOp.HTMLRecords() / "+errorID,
			"",
			"",
		)
	}

	context := map[string]string{
		"title":   "все з теґом '" + tag + "'",
		"header":  "Все з теґом '" + tag + "'",
		"content": htmlList,
	}
	return mustache.Render(htmlOp.htmlTemplate, context)
}

func (htmlOp *notebookHTML) HTMLTags(tagsStatMap tags.StatMap, options *crud.Options) (string, error) {

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

	return htmlTags, nil

}

func (htmlOp *notebookHTML) HTMLRecords(recordItems []records.Item, options *crud.Options) (string, error) {

	jsonBytes, err := json.Marshal(recordItems)

	return string(jsonBytes), err

}
