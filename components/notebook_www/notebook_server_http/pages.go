package notebook_server_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/selectors"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/data_exchange/components/tags"
	"github.com/pavlo67/tools/components/notebook_www"
	"github.com/pavlo67/tools/components/notebook_www/notebook_server_http/notebook_html"
	"github.com/pavlo67/tools/entities/records"
)

var PagesConfig = server_http.ConfigPages{
	ConfigCommon: server_http.ConfigCommon{
		Title:   "Notebook pages",
		Version: "0.0.1",
	},

	EndpointsPageSettled: server_http.EndpointsPageSettled{
		notebook_www.IntefaceKeyHTMLRoot:   {Path: "/", EndpointPage: rootPage},
		notebook_www.IntefaceKeyHTMLView:   {Path: "/view", EndpointPage: viewPage},
		notebook_www.IntefaceKeyHTMLCreate: {Path: "/create", EndpointPage: createPage},
		notebook_www.IntefaceKeyHTMLEdit:   {Path: "/edit", EndpointPage: editPage},
		notebook_www.IntefaceKeyHTMLSave:   {Path: "/save", EndpointPage: savePage},
		notebook_www.IntefaceKeyHTMLDelete: {Path: "/delete", EndpointPage: deletePage},
		notebook_www.IntefaceKeyHTMLTags:   {Path: "/tags", EndpointPage: tagsPage},
		notebook_www.IntefaceKeyHTMLTagged: {Path: "/tagged", EndpointPage: taggedPage},
		// notebook.IntefaceKeyHTMLList: {Path: "/list"},
	},
}

var rootPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(_ server_http.OperatorV2, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		tagsStatMap, err := recordsOp.Tags(nil, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.Tags()", req)
		}

		htmlIndex := notebookHTMLOp.HTMLIndex(identity)
		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, identity)

		fragments := server_http.CommonFragments(
			"вхід",
			"Вхід",
			"", "", htmlIndex,
			"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
		)

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var viewPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "GET",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := records.ID(params["record_id"])
		r, children, err := records.ReadWithChildren(recordsOp, id, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.ReadWithChildren()", req)
		}

		fragments, err := notebookHTMLOp.FragmentsView(r, children, "", identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var editPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "GET",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := records.ID(params["record_id"])

		r, err := recordsOp.Read(id, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.Read()", req)
		}

		fragments, err := notebookHTMLOp.FragmentsEdit(r, nil, "", identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при notebookHTMLOp.FragmentsEdit()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var createPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		fragments, err := notebookHTMLOp.FragmentsEdit(nil, nil, "", identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при notebookHTMLOp.FragmentsEdit()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var savePage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "POST",
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ErrorPage(http.StatusBadRequest, err, "при ioutil.ReadAll(req.Body)", req)
		}

		data, err := url.ParseQuery(string(body))
		if err != nil {
			return server_http.ErrorPage(http.StatusBadRequest, err, "при url.ParseQuery(body)", req)
		}

		r := notebook_html.RecordFromData(data)
		if r == nil {
			return server_http.ErrorPage(http.StatusBadRequest, fmt.Errorf("on notebook_html.RecordFromData(%#v): got nil", data), "при notebook_html.RecordFromData()", req)
		}

		r.ID, err = recordsOp.Save(*r, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.Save()", req)
		} else if r.ID == "" {
			return server_http.ErrorPage(0, fmt.Errorf("on recordsOp.Save(%#v, %#v): got nil", *r, identity), "при recordsOp.Save()", req)
		}

		r, children, err := records.ReadWithChildren(recordsOp, r.ID, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при ReadWithChildren()", req)
		}

		fragments, err := notebookHTMLOp.FragmentsView(r, children, "", identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var deletePage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "POST",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := records.ID(params["record_id"])

		err := recordsOp.Remove(id, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.Remove()", req)
		}

		fragments := server_http.CommonFragments(
			"запис вилучено",
			"Запис вилучено",
			"", "", "", "",
		)

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil
	},
}

var tagsPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		tagsStatMap, err := recordsOp.Tags(nil, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.Tags()", req)
		}

		htmlTags := notebookHTMLOp.HTMLTags(tagsStatMap, identity)
		//if err != nil {
		//	return ErrorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.HTMLTags()", req)
		//}

		fragments := server_http.CommonFragments(
			"теґи",
			"Теґи",
			"", "", "",
			"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
		)

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil

	},
}

var taggedPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "GET",
		PathParams: []string{"tag"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		tag := tags.Item(params["tag"])

		selectorTagged := selectors.Term{
			Key:    records.HasTag,
			Values: []string{tag},
		}

		rs, err := recordsOp.List(&selectorTagged, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при recordsOp.List()", req)
		}

		fragments, err := notebookHTMLOp.FragmentsListTagged(tag, rs, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: fragments,
		}, nil

	},
}
