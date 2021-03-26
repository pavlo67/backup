package files_server_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/selectors"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/data_exchange/components/tags"

	"github.com/pavlo67/tools/entities/files"

	"github.com/pavlo67/tools/components/files_www"
	"github.com/pavlo67/tools/components/files_www/files_server_http/files_html"
)

var PagesConfig = server_http.ConfigPages{
	ConfigCommon: server_http.ConfigCommon{
		Title:   "Files pages",
		Version: "0.0.1",
	},

	EndpointsPageSettled: server_http.EndpointsPageSettled{
		files_www.IntefaceKeyHTMLRoot:   {Path: "/", EndpointPage: rootPage},
		files_www.IntefaceKeyHTMLView:   {Path: "/view", EndpointPage: viewPage},
		files_www.IntefaceKeyHTMLCreate: {Path: "/create", EndpointPage: createPage},
		files_www.IntefaceKeyHTMLEdit:   {Path: "/edit", EndpointPage: editPage},
		files_www.IntefaceKeyHTMLSave:   {Path: "/save", EndpointPage: savePage},
		files_www.IntefaceKeyHTMLDelete: {Path: "/delete", EndpointPage: deletePage},
		// files.IntefaceKeyHTMLList: {Path: "/list"},
	},
}

var rootPage = server_http.EndpointPage{
	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(_ server_http.OperatorV2, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {

		htmlIndex := filesHTMLOp.HTMLIndex(identity)

		htmlPage, errRender := filesHTMLOp.CommonPage(
			"вхід",
			"Вхід",
			"", "", htmlIndex,
			"Розділи (теми) цієї бази даних: \n<p>",
		)
		if errRender != nil {
			return errorPage(0, filesHTMLOp, errRender, "при filesHTMLOp.CommonPage()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil
	},
}

var viewPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "GET",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := files.ID(params["record_id"])
		r, children, err := files.ReadWithChildren(filesOp, id, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.ReadWithChildren()", req)
		}

		htmlPage, err := filesHTMLOp.View(r, children, "", identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.View()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil
	},
}

var editPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "GET",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := files.ID(params["record_id"])

		r, err := filesOp.Read(id, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.Read()", req)
		}

		htmlPage, err := filesHTMLOp.Edit(r, nil, "", identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.Edit()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil
	},
}

var createPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		htmlPage, err := filesHTMLOp.Edit(nil, nil, "", identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.Edit()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
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
			return errorPage(http.StatusBadRequest, filesHTMLOp, err, "при ioutil.ReadAll(req.Body)", req)
		}

		data, err := url.ParseQuery(string(body))
		if err != nil {
			return errorPage(http.StatusBadRequest, filesHTMLOp, err, "при url.ParseQuery(body)", req)
		}

		r := files_html.RecordFromData(data)
		if r == nil {
			return errorPage(http.StatusBadRequest, filesHTMLOp, fmt.Errorf("on files_html.RecordFromData(%#v): got nil", data), "при files_html.RecordFromData()", req)
		}

		r.ID, err = filesOp.Save(*r, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.Save()", req)
		} else if r.ID == "" {
			return errorPage(0, filesHTMLOp, fmt.Errorf("on filesOp.Save(%#v, %#v): got nil", *r, identity), "при filesOp.Save()", req)
		}

		r, children, err := files.ReadWithChildren(filesOp, r.ID, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при ReadWithChildren()", req)
		}

		htmlPage, err := filesHTMLOp.View(r, children, "", identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.View()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil
	},
}

var deletePage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method:     "POST",
		PathParams: []string{"record_id"},
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		id := files.ID(params["record_id"])

		err := filesOp.Remove(id, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.Remove()", req)
		}

		htmlPage, errRender := filesHTMLOp.CommonPage(
			"запис вилучено",
			"Запис вилучено",
			"", "", "", "",
		)
		if errRender != nil {
			return errorPage(0, filesHTMLOp, errRender, "при filesHTMLOp.CommonPage()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil
	},
}

var tagsPage = server_http.EndpointPage{

	EndpointDescription: server_http.EndpointDescription{
		Method: "GET",
	},
	WorkerHTTPPage: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server_http.ResponsePage, error) {
		tagsStatMap, err := filesOp.Tags(nil, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.Tags()", req)
		}

		htmlTags := filesHTMLOp.HTMLTags(tagsStatMap, identity)
		//if err != nil {
		//	return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.HTMLTags()", req)
		//}

		htmlPage, errRender := filesHTMLOp.CommonPage(
			"теґи",
			"Теґи",
			"", "", "",
			"Розділи (теми) цієї бази даних: \n<p>"+htmlTags,
		)
		if errRender != nil {
			return errorPage(0, filesHTMLOp, errRender, "при filesHTMLOp.CommonPage()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
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
			Key:    files.HasTag,
			Values: []string{tag},
		}

		rs, err := filesOp.List(&selectorTagged, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesOp.List()", req)
		}

		htmlPage, err := filesHTMLOp.ListTagged(tag, rs, identity)
		if err != nil {
			return errorPage(0, filesHTMLOp, err, "при filesHTMLOp.View()", req)
		}

		return server_http.ResponsePage{
			Status:    http.StatusOK,
			Fragments: htmlPage,
		}, nil

	},
}

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
//		pathCorrected, err := filesOp.Save(bucketID, path, newFilePattern, data)
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
//		data, err := filesOp.Read(bucketID, path)
//		if err != nil || data == nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, data, req)
//	},
//}
//
//var removeEndpoint = server_http.Endpoint{
//	Method:     "DELETE",
//	PathParams: []string{"bucket_id", "path"},
//	WorkerHTTP: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
//		bucketID := files.BucketID(params["bucket_id"])
//		path := params["path"]
//
//		if err := filesOp.Remove(bucketID, path); err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, nil, req)
//	},
//}
//
//var listEndpoint = server_http.Endpoint{
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
//		filesInfo, err := filesOp.List(bucketID, path, depth)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, filesInfo, req)
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
//		fileInfo, err := filesOp.Stat(bucketID, path, depth)
//		if err != nil {
//			return serverOp.ResponseRESTError(0, err, req)
//		}
//
//		return serverOp.ResponseRESTOk(0, fileInfo, req)
//	},
//}
