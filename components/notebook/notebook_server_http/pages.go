package notebook_server_http

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/notebook/notebook_html"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

var Pages = server_http.Endpoints{
	rootPage,
	viewPage,
	editPage,
	savePage,
	//listPage,

	tagsPage,
	taggedPage,
}

var rootPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLRoot,
	Method:      "GET",
	WorkerHTTP: func(_ server_http.Operator, _ *http.Request, _ server_http.Params, options *crud.Options) (server.Response, error) {
		tagsStat, err := recordsOp.Tags(options)
		if err != nil {
			l.Error(err)
			return notebookHTMLOp.HTMLError(0, "На жаль, виникла помилка (при recordsOp.Tags(options))")
		}

		tagsStats := tagsStat.List(true)

		htmlTags, err := notebookHTMLOp.HTMLTags(tagsStats)
		if err != nil {
			l.Error(err)
			return notebookHTMLOp.HTMLError(0, "На жаль, виникла помилка (при notebookHTMLOp.HTMLTags(tagsStats))")
		}

		return notebookHTMLOp.HTMLRoot("Hello, World!", htmlTags)
	},
}

func htmlView(id records.ID, options *crud.Options) (server.Response, error) {
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

	message := notebookHTMLOp.HTMLMessage(errs)

	return notebookHTMLOp.HTMLView(r, children, message)
}

var viewPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		return htmlView(records.ID(params["record_id"]), options)
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

		message := notebookHTMLOp.HTMLMessage(errs)

		return notebookHTMLOp.HTMLEdit(r, nil, message)
	},
}

var savePage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLSave,
	Method:      "POST",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			l.Error(err)
		}

		//var data common.Map
		//if err = json.Unmarshal(dataBytes, &data); err != nil {
		//	l.Error(err)
		//}

		data, err := url.ParseQuery(string(body))
		if err != nil {
			l.Error(err)
		}

		r := notebook_html.RecordFromData(data)
		if r == nil {
			l.Errorf("no data??? %s", body)
			return notebookHTMLOp.HTMLPage("???", "??? no data", "", "ok!", ""), nil
		}

		l.Infof("$#v", r)

		r, err = recordsOp.Save(*r, options)
		if err != nil || r == nil {
			l.Errorf("??? %#v, %s, %#v", r, err, err)
			return notebookHTMLOp.HTMLPage("???", "???", "", "ok!", ""), nil
		}

		return htmlView(r.ID, options)
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

		htmlTags, err := notebookHTMLOp.HTMLTags(tagsStatList)
		errs = errs.Append(err)

		title := "нотатник: теґи"
		htmlHeader := "Теґи"

		return notebookHTMLOp.HTMLPage(title, htmlHeader, "", htmlTags, errs.Error()), nil
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

		htmlStr, err := notebookHTMLOp.HTMLList(tag, rs)
		errs = errs.Append(err)

		title := "нотатник: все з теґом '" + string(tag) + "'"
		htmlHeader := "Нотатник: все з теґом '" + string(tag) + "'"

		return notebookHTMLOp.HTMLPage(title, htmlHeader, "", htmlStr, errs.Error()), nil
	},
}

//selectorNoTag, err := recordsOp.HasNoTag()
//if err != nil {
//	l.Error(err)
//	return notebookHTMLOp.HTMLError(0, "На жаль, виникла помилка (при recordsOp.HasNoTag())")
//}
//
//optionsWithNoTag := options.WithSelector(selectorNoTag)
//
//rs, err := recordsOp.List(optionsWithNoTag)
//if err != nil {
//	l.Error(err)
//	return notebookHTMLOp.HTMLError(0, "На жаль, виникла помилка (при recordsOp.List(optionsWithNoTag))")
//}
//
