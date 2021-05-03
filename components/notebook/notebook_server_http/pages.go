package notebook_server_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/selectors"
	server_http_v2 "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/data/components/tags"

	"github.com/pavlo67/data/entities/records"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http/notebook_html"
)

const onPages = "on notebook_server_http.newNotebookPages()"

func newNotebookPages(prefix string, recordsOp records.Operator) (*server_http_v2.ConfigPages, error) {
	if recordsOp == nil {
		return nil, errors.New(onPages + ": no records.Operator")
	}

	pages := notebookPages{
		prefix:    prefix,
		recordsOp: recordsOp,
	}

	configPages := server_http_v2.ConfigPages{
		ConfigCommon: server_http.ConfigCommon{
			Title:   "Notebook pages",
			Version: "0.0.1",
			//Host:    "",
			//Port:    "",
			Prefix: prefix,
		},
		EndpointsPageSettled: server_http_v2.EndpointsPageSettled{
			notebook.IntefaceKeyHTMLRoot:   {Path: "/", EndpointPage: pages.root()},
			notebook.IntefaceKeyHTMLView:   {Path: "/view", EndpointPage: pages.view()},
			notebook.IntefaceKeyHTMLCreate: {Path: "/create", EndpointPage: pages.create()},
			notebook.IntefaceKeyHTMLEdit:   {Path: "/edit", EndpointPage: pages.edit()},
			notebook.IntefaceKeyHTMLSave:   {Path: "/save", EndpointPage: pages.save()},
			notebook.IntefaceKeyHTMLDelete: {Path: "/delete", EndpointPage: pages.delete()},
			notebook.IntefaceKeyHTMLTags:   {Path: "/tags", EndpointPage: pages.tags()},
			notebook.IntefaceKeyHTMLTagged: {Path: "/tagged", EndpointPage: pages.tagged()},
			// notebook.IntefaceKeyHTMLList: {Path: "/list"},
		},
	}

	var err error
	pages.notebookHTMLOp, err = notebook_html.New(configPages, l)
	if err != nil || pages.notebookHTMLOp == nil {
		return nil, fmt.Errorf(onPages+": on notebook_html.New() got %#v / %s", pages.notebookHTMLOp, err)
	}

	return &configPages, nil
}

type notebookPages struct {
	prefix         string
	recordsOp      records.Operator
	notebookHTMLOp *notebook_html.HTMLOp
}

func (pages *notebookPages) root() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method: "GET",
		},
		WorkerHTTPPage: func(_ server_http_v2.OperatorV2, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			tagsStatMap, err := pages.recordsOp.Tags(nil, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.Tags()", req)
			}

			htmlIndex := pages.notebookHTMLOp.HTMLIndex(identity)
			htmlTags := pages.notebookHTMLOp.HTMLTags(tagsStatMap, identity)

			fragments := server_http_v2.CommonFragments(
				"вхід",
				"Вхід",
				"", "", htmlIndex,
				"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
			)

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) view() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:     "GET",
			PathParams: []string{"record_id"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			id := records.ID(params["record_id"])
			r, children, err := records.ReadWithChildren(pages.recordsOp, id, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.ReadWithChildren()", req)
			}

			fragments, err := pages.notebookHTMLOp.FragmentsView(r, children, "", identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
			}

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) edit() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:     "GET",
			PathParams: []string{"record_id"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			id := records.ID(params["record_id"])

			r, err := pages.recordsOp.Read(id, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.Read()", req)
			}

			fragments, err := pages.notebookHTMLOp.FragmentsEdit(r, nil, "", identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при notebookHTMLOp.FragmentsEdit()", req)
			}

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) create() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method: "GET",
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			fragments, err := pages.notebookHTMLOp.FragmentsEdit(nil, nil, "", identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при notebookHTMLOp.FragmentsEdit()", req)
			}

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) save() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method: "POST",
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return server_http_v2.ErrorPage(http.StatusBadRequest, err, "при ioutil.ReadAll(req.Body)", req)
			}

			data, err := url.ParseQuery(string(body))
			if err != nil {
				return server_http_v2.ErrorPage(http.StatusBadRequest, err, "при url.ParseQuery(body)", req)
			}

			r := notebook_html.RecordFromData(data)
			if r == nil {
				return server_http_v2.ErrorPage(http.StatusBadRequest, fmt.Errorf("on notebook_html.RecordFromData(%#v): got nil", data), "при notebook_html.RecordFromData()", req)
			}

			r.ID, err = pages.recordsOp.Save(*r, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.Save()", req)
			} else if r.ID == "" {
				return server_http_v2.ErrorPage(0, fmt.Errorf("on recordsOp.Save(%#v, %#v): got nil", *r, identity), "при recordsOp.Save()", req)
			}

			r, children, err := records.ReadWithChildren(pages.recordsOp, r.ID, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при ReadWithChildren()", req)
			}

			fragments, err := pages.notebookHTMLOp.FragmentsView(r, children, "", identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
			}

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) delete() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:     "POST",
			PathParams: []string{"record_id"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			id := records.ID(params["record_id"])

			err := pages.recordsOp.Remove(id, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.Remove()", req)
			}

			fragments := server_http_v2.CommonFragments(
				"запис вилучено",
				"Запис вилучено",
				"", "", "", "",
			)

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil
		},
	}
}

func (pages *notebookPages) tags() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method: "GET",
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			tagsStatMap, err := pages.recordsOp.Tags(nil, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.Tags()", req)
			}

			htmlTags := pages.notebookHTMLOp.HTMLTags(tagsStatMap, identity)
			//if err != nil {
			//	return ErrorPage(0, notebookHTMLOp, err, "при notebookHTMLOp.HTMLTags()", req)
			//}

			fragments := server_http_v2.CommonFragments(
				"теґи",
				"Теґи",
				"", "", "",
				"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
			)

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil

		},
	}
}

func (pages *notebookPages) tagged() server_http_v2.EndpointPage {
	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:     "GET",
			PathParams: []string{"tag"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			tag := tags.Item(params["tag"])

			selectorTagged := selectors.Term{
				Key:    records.HasTag,
				Values: []string{tag},
			}

			rs, err := pages.recordsOp.List(&selectorTagged, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при recordsOp.List()", req)
			}

			fragments, err := pages.notebookHTMLOp.FragmentsListTagged(tag, rs, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при notebookHTMLOp.FragmentsView()", req)
			}

			return server_http_v2.ResponsePage{
				Status:    http.StatusOK,
				Fragments: fragments,
			}, nil

		},
	}
}
