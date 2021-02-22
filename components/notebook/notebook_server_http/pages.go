package notebook_server_http

import (
	"fmt"
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
	createPage,
	editPage,
	savePage,
	deletePage,

	tagsPage,
	taggedPage,
}

func errorPage(httpStatus int, notebookHTMLOp notebook_html.Operator, err error, publicDetails string, req *http.Request) (server.Response, error) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	htmlPage, errRender := notebookHTMLOp.CommonPage(
		"помилка",
		"",
		"",
		publicDetails,
		"",
		"",
	)

	var errs []interface{}

	if err != nil {
		errs = []interface{}{err}
	}
	if errRender != nil {
		errs = append(errs, errRender)
	}

	if len(errs) > 0 {
		if req != nil {
			err = errors.CommonError(append([]interface{}{fmt.Errorf("on %s %s", req.Method, req.URL)}, errs...)...)
		} else {
			err = errors.CommonError(errs...)
		}
	}

	return server.Response{
		Status:   http.StatusOK,
		Data:     []byte(htmlPage),
		MIMEType: "text/html; charset=utf-8",
	}, err
}

var rootPage = server_http.Endpoint{
	InternalKey: notebook.IntefaceKeyHTMLRoot,
	Method:      "GET",
	WorkerHTTP: func(_ server_http.Operator, req *http.Request, _ server_http.Params, options *crud.Options) (server.Response, error) {
		tagsStatMap, err := recordsOp.Tags(options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Tags()", req)
		}

		htmlIndex := notebookHTMLOp.HTMLIndex(options)
		//if err != nil {
		//	return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.HTMLIndex()", req)
		//}

		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, options)
		//if err != nil {
		//	return errorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.HTMLTags()", req)
		//}

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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])
		r, children, err := prepareRecord(id, options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при prepareRecord()", req)
		}

		htmlPage, err := notebookHTMLOp.View(r, children, "", options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Read()", req)
		}

		htmlPage, err := notebookHTMLOp.Edit(r, nil, "", options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		htmlPage, err := notebookHTMLOp.Edit(nil, nil, "", options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
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

		r, err = recordsOp.Save(*r, options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Save()", req)
		} else if r == nil {
			return errorPage(0, notebookHTMLOp, fmt.Errorf("on recordsOp.Save(%#v, %#v): got nil", *r, options), "при recordsOp.Save()", req)
		}

		r, children, err := prepareRecord(r.ID, options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при prepareRecord()", req)
		}

		htmlPage, err := notebookHTMLOp.View(r, children, "", options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		id := records.ID(params["record_id"])

		err := recordsOp.Remove(id, options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tagsStatMap, err := recordsOp.Tags(options)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.Tags()", req)
		}

		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, options)
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		tag := tags.Item(params["tag"])

		selectorTagged, err := recordsOp.HasTag(tag)
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.HasTag()", req)
		}

		rs, err := recordsOp.List(options.WithSelector(selectorTagged))
		if err != nil {
			return errorPage(0, notebookHTMLOp, err, "при recordsOp.List()", req)
		}

		htmlPage, err := notebookHTMLOp.ListTagged(tag, rs, options)
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
