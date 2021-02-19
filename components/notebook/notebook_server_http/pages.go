package notebook_server_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/crud"
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
	createPage,
	editPage,
	savePage,
	deletePage,

	tagsPage,
	taggedPage,
}

var rootPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLRoot,
	Method:      "GET",
	WorkerHTTP: func(_ server_http.Operator, req *http.Request, _ server_http.Params, options *crud.Options) (server.Response, error) {
		tagsStatMap, err := recordsOp.Tags(options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.Tags()", req)
		}

		return notebookHTMLOp.HTMLRoot("Hello, World!", tagsStatMap)
	},
}

var viewPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])
		r, children, err := prepareRecord(id, options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при prepareRecord()", req)
		}

		return notebookHTMLOp.HTMLView(r, children, "")
	},
}

var editPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.Read()", req)
		}

		return notebookHTMLOp.HTMLEdit(r, nil, "")
	},
}

var createPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.Read()", req)
		}

		return notebookHTMLOp.HTMLEdit(r, nil, "")
	},
}

var savePage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLSave,
	Method:      "POST",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return notebookHTMLOp.HTMLError(http.StatusBadRequest, err, "при ioutil.ReadAll(req.Body)", req)
		}

		data, err := url.ParseQuery(string(body))
		if err != nil {
			return notebookHTMLOp.HTMLError(http.StatusBadRequest, err, "при url.ParseQuery(body)", req)
		}

		r := notebook_html.RecordFromData(data)
		if r == nil {
			return notebookHTMLOp.HTMLError(http.StatusBadRequest, fmt.Errorf("from %#v", data), "при notebook_html.RecordFromData()", req)
		}

		r, err = recordsOp.Save(*r, options)
		if err != nil || r == nil {
			return notebookHTMLOp.HTMLError(0, fmt.Errorf("got %#v, %s", r, err), "при recordsOp.Save()", req)
		}

		r, children, err := prepareRecord(r.ID, options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при prepareRecord()", req)
		}

		return notebookHTMLOp.HTMLView(r, children, "")
	},
}

var deletePage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLSave,
	Method:      "POST",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		err := recordsOp.Remove(id, options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.Remove()", req)
		}

		return notebookHTMLOp.HTMLRoot("запис вилучено!", nil)
	},
}

var tagsPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTags,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {

		tagsStatMap, err := recordsOp.Tags(options)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.Tags()", req)
		}

		return notebookHTMLOp.HTMLTags(tagsStatMap)
	},
}

var taggedPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTagged,
	Method:      "GET",
	PathParams:  []string{"tag"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tag := tags.Item(params["tag"])

		selectorTagged, err := recordsOp.HasTag(tag)
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.HasTag()", req)
		}

		rs, err := recordsOp.List(options.WithSelector(selectorTagged))
		if err != nil {
			return notebookHTMLOp.HTMLError(0, err, "при recordsOp.List()", req)
		}

		return notebookHTMLOp.HTMLTagged(tag, rs)
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
