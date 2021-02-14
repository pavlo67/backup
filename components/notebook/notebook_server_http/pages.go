package notebook_server_http

import (
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

var Endpoints = server_http.Endpoints{
	rootEndpoint,
	viewEndpoint,
	editEndpoint,
	tagsEndpoint,
	taggedEndpoint,
}

var rootEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLRoot,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {
		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte("мама мила раму!!!"),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var viewEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			l.Error(err)
		}

		var children []records.Item

		selector, err := recordsOp.HasParent(id)
		if err != nil {
			l.Error(err)
		} else {
			options = options.WithSelector(selector)
			children, err = recordsOp.List(options)
			if err != nil {
				l.Error(err)
			}
		}

		htmlStr, err := recordsHTMLOp.HTMLView(r, children)
		if err != nil {
			l.Error(err)
		}

		return ResponseHTMLOk(0, htmlStr), nil
	},
}

var editEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			l.Error(err)
		}

		htmlStr := fmt.Sprintf("edit form for %s --> %#v", id, r)
		return ResponseHTMLOk(0, htmlStr), nil
	},
}

var tagsEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTags,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tagsStat, err := recordsOp.Tags(options)
		if err != nil {
			l.Error(err)
		}

		tagsStatList := tagsStat.List(true)

		htmlTags, err := tagsHTMLOp.HTMLTags(tagsStatList)
		if err != nil {
			l.Error(err)
		}

		return ResponseHTMLOk(0, htmlTags), nil
	},
}

var taggedEndpoint = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTagged,
	Method:      "GET",
	PathParams:  []string{"tag"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tag := tags.Item(params["tag"])

		selectorTagged, err := recordsOp.HasTag(tag)
		if err != nil {
			l.Error(err)
		}

		optionsWithTag := options.WithSelector(selectorTagged)

		rs, err := recordsOp.List(optionsWithTag)
		if err != nil {
			l.Error(err)
		}

		htmlStr, err := recordsHTMLOp.HTMLTagged(tag, rs)
		if err != nil {
			l.Error(err)
		}

		return ResponseHTMLOk(0, htmlStr), nil
	},
}
