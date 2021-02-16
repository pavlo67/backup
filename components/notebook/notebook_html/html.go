package notebook_html

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/tools/components/ns"

	"github.com/cbroglie/mustache"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
	server_http2 "github.com/pavlo67/tools/common/server/server_http"

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

	method, urlStr, err := server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLView, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLList, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLList, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLList, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLEdit, []string{"id"}, false)
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

const onPrepare = "on notebookHTML.Prepare(): "

func (htmlOp *notebookHTML) Prepare(key Key, template string, params common.Map) error {
	return nil
}

var dataFields = []views_html.Field{
	{"id", "", "hidden", nil, nil},
	{"issued_id", "", "hidden", nil, nil},
	{"data_type", "", "hidden", nil, nil},
	// {"visibility", "тип", "select", nil, "ut"},
	// {"history_key", "", "hidden", nil, nil},

	{"title", "заголовок", "", nil, nil},
	{"summary", "коротко про", "", nil, nil},
	{"content_data", "опис", "", common.Map{"format": "35"}, nil},
	{"tags", "теми, розділи", "tag-it", nil, nil},
	// {"data_subtype", "", "hidden", "", ""},
	// {"embedded", "", "hidden", "", ""},
	// {"files", "завантажити файл", "file", "", "ut"},

	{"created_at", "створено", "view", common.Map{"format": "datetime"}, map[string]string{"class": "not_empty"}},
	{"updated_at", "востаннє відредаґовано", "view", common.Map{"format": "datetime"}, map[string]string{"class": "not_empty"}},
}

func value(data map[string][]string, key string) string {
	v := data[key]
	if len(v) == 1 {
		return v[0]
	} else if len(v) > 1 {
		return strings.Join(v, " ")
	}

	return ""
}

func RecordFromData(data map[string][]string) *records.Item {
	if data == nil {
		return nil
	}

	var tagItems []tags.Item
	for _, t := range strings.Split(value(data, "tags"), ";") {
		tagItems = append(tagItems, tags.Item(strings.TrimSpace(t)))
	}

	r := records.Item{
		ID:       records.ID(value(data, "id")),
		IssuedID: ns.ID(value(data, "issued_id")),
		Content: records.Content{
			Title:   value(data, "title"),
			Summary: value(data, "summary"),
			// TypeKey:  "",
			Data: value(data, "content_data"),
			// Embedded: nil,
			Tags: tagItems,
		},
		//OwnerID:   "",
		//ViewerID:  "",
	}

	return &r
}

func DataFromRecord(r *records.Item) map[string]string {
	if r == nil {
		return nil
	}

	var updatedAt string
	if r.UpdatedAt != nil {
		updatedAt = r.UpdatedAt.Format("02.01.2006 15:04:05")
	}

	data := map[string]string{
		"id":           string(r.ID),
		"issued_id":    string(r.IssuedID),
		"data_type":    "record", // TODO!!!
		"title":        r.Content.Title,
		"summary":      r.Content.Summary,
		"content_data": r.Content.Data,
		"tags":         strings.Join(r.Content.Tags, "; "),
		// "embedded": r.Content.Embedded,
		"created_at": r.CreatedAt.Format("02.01.2006 15:04:05"),
		"updated_at": updatedAt,
	}

	//linksList, err := json.Marshal(r.Links)
	//if err != nil {
	//	return nil, nil, errors.Wrapf(err, "can't marshal object.tags: %#v for object.id: %s", r.Links, r.ID)
	//}
	//DataFromRecord["links"] = string(linksList)
	//
	//tags := ""
	//filesList := []interfaces.Link{}
	//for _, l := range r.Links {
	//	switch l.Type {
	//
	//	case links.TypeTag:
	//		tags += l.Name + "; "
	//
	//	case files.LinkType:
	//		filesList = append(filesList, l)
	//	}
	//}
	//DataFromRecord["tags"] = tags
	//if len(filesList) > 0 {
	//	files, err := json.Marshal(filesList)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	DataFromRecord["files"] = string(files)
	//}
	//
	//if r.UpdatedAt != nil {
	//	DataFromRecord["updated_at"] = r.UpdatedAt.Format("02.01.2006 15:04:05")
	//}

	return data

}

var createFields = append(dataFields, views_html.Field{"create", "зберегти запис", "button", nil, map[string]string{"class": "ut"}})

// TODO!!! look at https://github.com/kataras/blocks

func (htmlOp *notebookHTML) HTMLMessage(errs errors.Error) string {
	return errs.Error()
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

		method, urlStr, err := server_http2.EP(htmlOp.pagesConfig, notebook.IntefaceKeyHTMLList, []string{tag}, false)
		if err != nil || urlStr == "" {
			l.Errorf("can't server_http.EP(%#v, notebook.IntefaceKeyHTMLTags,nil,false), got %s, %s, %s", htmlOp.pagesConfig, method, urlStr, err)
			continue
		}

		htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	}

	htmlTags += "</ul>\n\n"

	return htmlTags, nil
}
