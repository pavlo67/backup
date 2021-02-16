package notebook_server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
	"github.com/pavlo67/tools/components/views/views_html"
)

var Pages = server_http.Endpoints{
	rootPage,
	viewPage,
	editPage,
	tagsPage,
	taggedPage,
}

var rootPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLRoot,
	Method:      "GET",
	WorkerHTTP: func(_ server_http.Operator, _ *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {
		return HTMLPage("нотатник", "Нотатник", "", "!!!", ""), nil
	},
}

var dataFields = []views_html.Field{
	{"visibility", "тип", "select", "", "ut"},
	{"title", "заголовок", "", "", ""},
	{"content", "опис", "", "35", ""},
	{"tags", "теми, розділи", "tag-it", "", ""},
	{"updated_at", "востаннє відредаґовано", "view", "datetime", "not_empty"},
	{"files", "Завантажити файл", "file", "", "ut"},
	{"id", "", "hidden", "", ""},
	{"genus", "", "hidden", "", ""},
}

var createFields = append(dataFields, views_html.Field{"create", "Зберегти запис", "button", "", "ut"})
var updateFields = append(dataFields, views_html.Field{"update", "Зберегти зміни", "button", "", "ut"})

var viewPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		errs := errors.CommonError()

		r, err := recordsOp.Read(id, options)
		errs = errs.Append(err)

		var children []records.Item

		selector, err := recordsOp.HasParent(id)
		if err != nil {
			errs = errs.Append(err)
		} else {
			options = options.WithSelector(selector)
			children, err = recordsOp.List(options)
			if err != nil {
				l.Error(err)
			}
		}

		title := "нотатник: " + r.Content.Title
		htmlHeader := r.Content.Title

		htmlStr, err := recordsHTMLOp.HTMLView(r, children)
		errs = errs.Append(err)

		return HTMLPage(title, htmlHeader, "", htmlStr, errs.Error()), nil

	},
}

var editPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		errs := errors.CommonError()

		r, err := recordsOp.Read(id, options)
		errs = errs.Append(err)

		title := "нотатник: " + r.Content.Title
		htmlHeader := r.Content.Title

		//htmlStr := fmt.Sprintf("edit form for %s --> %#v", id, r)
		htmlStr := HTMLEditTable(dataFields, "nb_edit_", nil, nil)

		return HTMLPage(title, htmlHeader, "", htmlStr, errs.Error()), nil
	},
}

var tagsPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTags,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {

		errs := errors.CommonError()

		tagsStat, err := recordsOp.Tags(options)
		errs = errs.Append(err)

		tagsStatList := tagsStat.List(true)

		htmlTags, err := tagsHTMLOp.HTMLTags(tagsStatList)
		errs = errs.Append(err)

		title := "нотатник: теґи"
		htmlHeader := "Теґи"

		return HTMLPage(title, htmlHeader, "", htmlTags, errs.Error()), nil
	},
}

var taggedPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTagged,
	Method:      "GET",
	PathParams:  []string{"tag"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tag := tags.Item(params["tag"])

		errs := errors.CommonError()

		selectorTagged, err := recordsOp.HasTag(tag)
		errs = errs.Append(err)

		optionsWithTag := options.WithSelector(selectorTagged)

		rs, err := recordsOp.List(optionsWithTag)
		errs = errs.Append(err)

		htmlStr, err := recordsHTMLOp.HTMLTagged(tag, rs)
		errs = errs.Append(err)

		title := "нотатник: все з теґом '" + string(tag) + "'"
		htmlHeader := "Нотатник: все з теґом '" + string(tag) + "'"

		return HTMLPage(title, htmlHeader, "", htmlStr, errs.Error()), nil
	},
}
