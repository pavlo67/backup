package catalogue_server_http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/server/server_http"

	server_http_v2 "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/data/entities/items"

	"github.com/pavlo67/tools/components/catalogue"
)

const onPages = "on catalogue_server_http.newNotebookPages()"

func newCataloguePages(prefix string, itemsOp items.Operator) (*server_http_v2.ConfigPages, error) {
	if itemsOp == nil {
		return nil, errors.New(onPages + ": no items.Operator")
	}

	pages := cataloguePages{
		prefix:  prefix,
		itemsOp: itemsOp,
	}

	configPages := server_http_v2.ConfigPages{
		ConfigCommon: server_http.ConfigCommon{
			Title:   "Catalogue pages",
			Version: "0.0.1",
			//Host:    "",
			//Port:    "",
			Prefix: prefix,
		},
		EndpointsPageSettled: server_http_v2.EndpointsPageSettled{
			catalogue.IntefaceKeyHTMLList:   {Path: "/list", EndpointPage: pages.list()},
			catalogue.IntefaceKeyHTMLDelete: {Path: "/delete", EndpointPage: pages.delete()},

			//files_www.IntefaceKeyHTMLView:   {Path: "/view", EndpointPage: viewPage},
			//files_www.IntefaceKeyHTMLCreate: {Path: "/create", EndpointPage: createPage},
			//files_www.IntefaceKeyHTMLEdit:   {Path: "/edit", EndpointPage: editPage},
			//files_www.IntefaceKeyHTMLSave:   {Path: "/save", EndpointPage: savePage},
		},
	}

	var err error
	pages.catalogueHTMLOp, err = New(configPages, prefix)
	if err != nil || pages.catalogueHTMLOp == nil {
		return nil, fmt.Errorf(onPages+": on catalogue_html.New() got %#v / %s", pages.catalogueHTMLOp, err)
	}

	return &configPages, nil
}

type cataloguePages struct {
	prefix          string
	itemsOp         items.Operator
	catalogueHTMLOp *catalogueHTML
}

func (pages *cataloguePages) list() server_http_v2.EndpointPage {

	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:      "GET",
			PathParams:  []string{"*path"},
			QueryParams: []string{"depth"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			path := params["path"]
			depth := 0
			//depth, err := strconv.Atoi(params["depth"])
			//if err != nil {
			//	errors.Wrapf(err, "can't read depth (%s)", params["depth"])
			//	return server_http_v2.ErrorPage(0, err, "при itemsOp.List()", req)
			//}

			filesItems, err := pages.itemsOp.List(path, depth, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при itemsOp.List()", req)
			}

			htmlPage, err := pages.catalogueHTMLOp.FragmentsList(path, filesItems, path, identity)
			if err != nil {
				return server_http_v2.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
			}

			return server_http_v2.ResponsePage{Status: http.StatusOK, Fragments: htmlPage}, nil
		},
	}
}

func (pages *cataloguePages) delete() server_http_v2.EndpointPage {

	return server_http_v2.EndpointPage{

		EndpointDescription: server_http.EndpointDescription{
			Method:     "POST",
			PathParams: []string{"*path"},
		},
		WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
			path := params["path"]
			if err := pages.itemsOp.Remove(path, identity); err != nil {
				return server_http_v2.ErrorPage(0, err, "при itemsOp.Remove()", req)
			}

			htmlPage := server_http_v2.CommonFragments(
				"файл вилучено: "+path,
				"Файл вилучено: "+path,
				"", "", "", "",
			)

			return server_http_v2.ResponsePage{Status: http.StatusOK, Fragments: htmlPage}, nil
		},
	}
}

//var viewPage = server_http_v2.EndpointPage{
//
//	EndpointDescription: server_http_v2.EndpointDescription{
//		Method:     "GET",
//		PathParams: []string{"record_id"},
//	},
//	WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
//		r, err := itemsOp.Read(path)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при itemsOp.ReadWithChildren()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsView(r, "", identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
//		}
//
//		return server_http_v2.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}

//var editPage = server_http_v2.EndpointPage{
//
//	EndpointDescription: server_http_v2.EndpointDescription{
//		Method:     "GET",
//		PathParams: []string{"record_id"},
//	},
//	WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
//
//		r, err := itemsOp.Read(id, identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при itemsOp.Read()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsEdit(r,  "", identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при filesHTMLOp.FragmentsEdit()", req)
//		}
//
//		return server_http_v2.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}
//
//var createPage = server_http_v2.EndpointPage{
//
//	EndpointDescription: server_http_v2.EndpointDescription{
//		Method: "GET",
//	},
//	WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
//		htmlPage, err := filesHTMLOp.FragmentsEdit(nil,  "", identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при filesHTMLOp.FragmentsEdit()", req)
//		}
//
//		return server_http_v2.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}
//
//var savePage = server_http_v2.EndpointPage{
//
//	EndpointDescription: server_http_v2.EndpointDescription{
//		Method: "POST",
//	},
//	WorkerHTTPPage: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server_http_v2.ResponsePage, error) {
//		body, err := ioutil.ReadAll(req.Body)
//		if err != nil {
//			return server_http_v2.ErrorPage(http.StatusBadRequest, err, "при ioutil.ReadAll(req.Body)", req)
//		}
//
//		data, err := url.ParseQuery(string(body))
//		if err != nil {
//			return server_http_v2.ErrorPage(http.StatusBadRequest, err, "при url.ParseQuery(body)", req)
//		}
//
//		r := files_html.RecordFromData(data)
//		if r == nil {
//			return server_http_v2.ErrorPage(http.StatusBadRequest, fmt.Errorf("on files_html.RecordFromData(%#v): got nil", data), "при files_html.RecordFromData()", req)
//		}
//
//		r.ID, err = itemsOp.Save(*r, identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при itemsOp.Save()", req)
//		} else if r.ID == "" {
//			return server_http_v2.ErrorPage(0, fmt.Errorf("on itemsOp.Save(%#v, %#v): got nil", *r, identity), "при itemsOp.Save()", req)
//		}
//
//		r, children, err := files.ReadWithChildren(itemsOp, r.ID, identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при ReadWithChildren()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsView(r, children, "", identity)
//		if err != nil {
//			return server_http_v2.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
//		}
//
//		return server_http_v2.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}

//var saveEndpoint = server_http_v2.EndpointREST{
//	Method:     "POST",
//	PathParams: []string{"bucket_id", "path", "new_file_pattern"},
//	WorkerHTTP: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//		newFilePattern := params["new_file_pattern"]
//
//		data, err := ioutil.ReadAll(req.Body)
//		if err != nil {
//			return serverOp.ResponseRESTError(http.StatusBadRequest, errors.CommonError(err, "reading body"), req)
//		}
//
//		pathCorrected, err := itemsOp.Save(bucketID, path, newFilePattern, data)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, pathCorrected, req)
//	},
//}
//
//var readEndpoint = server_http_v2.EndpointREST{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path"},
//	WorkerHTTP: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//
//		data, err := itemsOp.Read(bucketID, path)
//		if err != nil || data == nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, data, req)
//	},
//}
//
//var statEndpoint = server_http_v2.EndpointREST{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path", "depth"},
//	WorkerHTTP: func(serverOp server_http_v2.OperatorV2, req *http.Request, params server_http_v2.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//		depth, err := strconv.Atoi(params["depth"])
//		if err != nil {
//			return serverOp.ResponseRESTError(0, errors.Wrapf(err, "can't read depth (%s)", params["depth"]), req)
//		}
//
//		fileInfo, err := itemsOp.Stat(bucketID, path, depth)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, fileInfo, req)
//	},
//}
