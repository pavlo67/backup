package files_server_http

//var saveEndpoint = server_http.EndpointREST{
//	Method:     "POST",
//	PathParams: []string{"bucket_id", "path", "new_file_pattern"},
//	WorkerHTTPREST: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
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
//var readEndpoint = server_http.EndpointREST{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path"},
//	WorkerHTTPREST: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
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
//var removeEndpoint = server_http.EndpointREST{
//	Method:     "DELETE",
//	PathParams: []string{"bucket_id", "path"},
//	WorkerHTTPREST: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
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
//var listEndpoint = server_http.EndpointREST{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path", "depth"},
//	WorkerHTTPREST: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
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
//var statEndpoint = server_http.EndpointREST{
//	Method:     "GET",
//	PathParams: []string{"bucket_id", "path", "depth"},
//	WorkerHTTPREST: func(serverOp server_http.OperatorV2, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
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
