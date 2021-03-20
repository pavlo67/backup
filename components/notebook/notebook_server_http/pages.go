package notebook_server_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/data_exchange/components/tags"
	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http/notebook_html"
	"github.com/pavlo67/tools/entities/records"
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
	WorkerHTTP: func(_ server_http.Operator, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		tagsStatMap, err := recordsOp.Tags(nil, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Tags()", req)
		}

		htmlIndex := notebookHTMLOp.HTMLIndex(identity)
		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, identity)

		htmlPage, errRender := notebookHTMLOp.CommonPage(
			"вхід",
			"Вхід",
			"", "", htmlIndex,
			"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
		)
		if errRender != nil {
			return errorPage(0, notebookHTMLOp, errRender, "при notebookHTMLOp.CommonPage()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var viewPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLView,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		id := records.ID(params["record_id"])
		r, children, err := records.ReadWithChildren(recordsOp, id, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.ReadWithChildren()", req)
		}

		htmlPage, err := notebookHTMLOp.View(r, children, "", identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.View()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var editPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLEdit,
	Method:      "GET",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Read()", req)
		}

		htmlPage, err := notebookHTMLOp.Edit(r, nil, "", identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.Edit()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var createPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLCreate,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		htmlPage, err := notebookHTMLOp.Edit(nil, nil, "", identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.Edit()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var savePage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLSave,
	Method:      "POST",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return errorPage(http.StatusBadRequest, notebookHTMLOp, err, "при ioutil.ReadAll(req.Body)", req)
		}

		data, err := url.ParseQuery(string(body))
		if err != nil {
			return errorPage(http.StatusBadRequest, notebookHTMLOp, err, "при url.ParseQuery(body)", req)
		}

		r := notebook_html.RecordFromData(data)
		if r == nil {
			return errorPage(http.StatusBadRequest, notebookHTMLOp, fmt.Errorf("on notebook_html.RecordFromData(%#v): got nil", data), "при notebook_html.RecordFromData()", req)
		}

		r.ID, err = recordsOp.Save(*r, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Save()", req)
		} else if r.ID == "" {
			return errorPage(0, notebookHTMLOp, fmt.Errorf("on recordsOp.Save(%#v, %#v): got nil", *r, identity), "при recordsOp.Save()", req)
		}

		r, children, err := records.ReadWithChildren(recordsOp, r.ID, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при ReadWithChildren()", req)
		}

		htmlPage, err := notebookHTMLOp.View(r, children, "", identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.View()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var deletePage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLDelete,
	Method:      "POST",
	PathParams:  []string{"record_id"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		id := records.ID(params["record_id"])

		err := recordsOp.Remove(id, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Remove()", req)
		}

		htmlPage, errRender := notebookHTMLOp.CommonPage(
			"запис вилучено",
			"Запис вилучено",
			"", "", "", "",
		)
		if errRender != nil {
			return errorPage(0, notebookHTMLOp, errRender, "при notebookHTMLOp.CommonPage()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil
	},
}

var tagsPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTags,
	Method:      "GET",
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		tagsStatMap, err := recordsOp.Tags(nil, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Tags()", req)
		}

		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, identity)
		//if err != nil {
		//	return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.HTMLTags()", req)
		//}

		htmlPage, errRender := notebookHTMLOp.CommonPage(
			"теґи",
			"Теґи",
			"", "", "",
			"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
		)
		if errRender != nil {
			return errorPage(0, notebookHTMLOp, errRender, "при notebookHTMLOp.CommonPage()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil

	},
}

var taggedPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLTagged,
	Method:      "GET",
	PathParams:  []string{"tag"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		tag := tags.Item(params["tag"])

		selectorTagged := selectors.Term{
			Key:    records.HasTag,
			Values: []string{tag},
		}

		rs, err := recordsOp.List(&selectorTagged, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.List()", req)
		}

		htmlPage, err := notebookHTMLOp.ListTagged(tag, rs, identity)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.View()", req)
		}

		return server.Response{
			Status:   http.StatusOK,
			Data:     []byte(htmlPage),
			MIMEType: "text/html; charset=utf-8",
		}, nil

	},
}
