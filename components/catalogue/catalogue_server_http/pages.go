package catalogue_server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/components/catalogue"
)

var PagesConfig = server_http.ConfigPages{
	ConfigCommon: server_http.ConfigCommon{
		Title:   "Files pages",
		Version: "0.0.1",
	},

	EndpointsPageSettled: server_http.EndpointsPageSettled{
		catalogue.IntefaceKeyHTMLList:   {Path: "/list", EndpointPage: listPage},
		catalogue.IntefaceKeyHTMLDelete: {Path: "/delete", EndpointPage: deletePage},

		//files_www.IntefaceKeyHTMLView:   {Path: "/view", EndpointPage: viewPage},
		//files_www.IntefaceKeyHTMLCreate: {Path: "/create", EndpointPage: createPage},
		//files_www.IntefaceKeyHTMLEdit:   {Path: "/edit", EndpointPage: editPage},
		//files_www.IntefaceKeyHTMLSave:   {Path: "/save", EndpointPage: savePage},
	},
}

var listPage = server_http.EndpointPage{
	EndpointDescription: server_http.EndpointDescription{
		Method:      "GET",
		PathParams:  []string{"*path"},
		QueryParams: []string{"depth"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		path := params["path"]
		depth := 0
		//depth, err := strconv.Atoi(params["depth"])
		//if err != nil {
		//	errors.Wrapf(err, "can't read depth (%s)", params["depth"])
		//	return server_http.ErrorPage(0, err, "при catalogueOp.List()", req)
		//}

		filesItems, err := catalogueOp.List(path, depth, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при catalogueOp.List()", req)
		}

		htmlPage, err := filesHTMLOp.FragmentsList(path, filesItems, path, identity)
		if err != nil {
			return server_http.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
		}

		return server_http.ResponsePage{Status: http.StatusOK, Fragments: htmlPage}, nil
	},
}

var deletePage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "POST",
		PathParams: []string{"*path"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		path := params["path"]
		if err := catalogueOp.Remove(path, identity); err != nil {
			return server_http.ErrorPage(0, err, "при catalogueOp.Remove()", req)
		}

		htmlPage := server_http.CommonFragments(
			"файл вилучено: "+path,
			"Файл вилучено: "+path,
			"", "", "", "",
		)

		return server_http.ResponsePage{Status: http.StatusOK, Fragments: htmlPage}, nil
	},
}

//var viewPage = server_http.EndpointPage{
//
//	EndpointDescription: server_http.EndpointDescription{
//		Method:     "GET",
//		PathParams: []string{"record_id"},
//	},
//	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
//		r, err := catalogueOp.Read(path)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при catalogueOp.ReadWithChildren()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsView(r, "", identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
//		}
//
//		return server_http.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}

//var editPage = server_http.EndpointPage{
//
//	EndpointDescription: server_http.EndpointDescription{
//		Method:     "GET",
//		PathParams: []string{"record_id"},
//	},
//	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
//
//		r, err := catalogueOp.Read(id, identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при catalogueOp.Read()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsEdit(r,  "", identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при filesHTMLOp.FragmentsEdit()", req)
//		}
//
//		return server_http.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}
//
//var createPage = server_http.EndpointPage{
//
//	EndpointDescription: server_http.EndpointDescription{
//		Method: "GET",
//	},
//	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
//		htmlPage, err := filesHTMLOp.FragmentsEdit(nil,  "", identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при filesHTMLOp.FragmentsEdit()", req)
//		}
//
//		return server_http.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}
//
//var savePage = server_http.EndpointPage{
//
//	EndpointDescription: server_http.EndpointDescription{
//		Method: "POST",
//	},
//	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
//		body, err := ioutil.ReadAll(req.Body)
//		if err != nil {
//			return server_http.ErrorPage(http.StatusBadRequest, err, "при ioutil.ReadAll(req.Body)", req)
//		}
//
//		data, err := url.ParseQuery(string(body))
//		if err != nil {
//			return server_http.ErrorPage(http.StatusBadRequest, err, "при url.ParseQuery(body)", req)
//		}
//
//		r := files_html.RecordFromData(data)
//		if r == nil {
//			return server_http.ErrorPage(http.StatusBadRequest, fmt.Errorf("on files_html.RecordFromData(%#v): got nil", data), "при files_html.RecordFromData()", req)
//		}
//
//		r.ID, err = catalogueOp.Save(*r, identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при catalogueOp.Save()", req)
//		} else if r.ID == "" {
//			return server_http.ErrorPage(0, fmt.Errorf("on catalogueOp.Save(%#v, %#v): got nil", *r, identity), "при catalogueOp.Save()", req)
//		}
//
//		r, children, err := files.ReadWithChildren(catalogueOp, r.ID, identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при ReadWithChildren()", req)
//		}
//
//		htmlPage, err := filesHTMLOp.FragmentsView(r, children, "", identity)
//		if err != nil {
//			return server_http.ErrorPage(0, err, "при filesHTMLOp.FragmentsView()", req)
//		}
//
//		return server_http.ResponsePage{
//			Status:    http.StatusOK,
//			Fragments: htmlPage,
//		}, nil
//	},
//}

//var saveEndpoint = server_http.Endpoint{
//	Method:     "POST",
//	PathParams: []string{"bucket_id", "path", "new_file_pattern"},
//	WorkerHTTP: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//		newFilePattern := params["new_file_pattern"]
//
//		data, err := ioutil.ReadAll(req.Body)
//		if err != nil {
//			return serverOp.ResponseRESTError(http.StatusBadRequest, errors.CommonError(err, "reading body"), req)
//		}
//
//		pathCorrected, err := catalogueOp.Save(bucketID, path, newFilePattern, data)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, pathCorrected, req)
//	},
//}
//
//var readEndpoint = server_http.Endpoint{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path"},
//	WorkerHTTP: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//
//		data, err := catalogueOp.Read(bucketID, path)
//		if err != nil || data == nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, data, req)
//	},
//}
//
//var statEndpoint = server_http.Endpoint{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path", "depth"},
//	WorkerHTTP: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//		depth, err := strconv.Atoi(params["depth"])
//		if err != nil {
//			return serverOp.ResponseRESTError(0, errors.Wrapf(err, "can't read depth (%s)", params["depth"]), req)
//		}
//
//		fileInfo, err := catalogueOp.Stat(bucketID, path, depth)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, fileInfo, req)
//	},
//}
