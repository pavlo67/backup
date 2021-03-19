package files_http

//var _ files.Operator = &filesHTTP{}
//
//type filesHTTP struct {
//	host      string
//	endpoints server_http.Endpoints
//}
//
//const onNew = "on filesHTTP.New()"
//
//func New(access config.Access, prefix string, endpoints server_http.Endpoints) (files.Operator, error) {
//	filesOp := filesHTTP{
//		host:      access.Host,
//		endpoints: server_http.Endpoints{},
//	}
//
//	if access.Port > 0 {
//		filesOp.host += ":" + strconv.Itoa(access.Port)
//	}
//	filesOp.host += prefix
//
//	return &filesOp, nil
//}

//func (filesOp *filesHTTP) Save(bucketID files.BucketID, path, newFilePattern string, data []byte, identity *auth.Identity) (string, error) {
//	if path = strings.TrimSpace(path); path == "" {
//		path = "."
//	}
//	if newFilePattern = strings.TrimSpace(newFilePattern); newFilePattern == "" {
//		newFilePattern = "."
//	}
//
//	ep := filesOp.endpoints[files.EPSave]
//	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + newFilePattern
//
//	var correctedPath string
//	return correctedPath, server_http.Request(serverURL, ep, nil, &correctedPath, options.GetIdentity(), filesOp.logfile)
//}
//
//func (filesOp *filesHTTP) Read(bucketID files.BucketID, path string, identity *auth.Identity) ([]byte, error) {
//	ep := filesOp.endpoints[files.EPRead]
//	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path
//
//	var data []byte
//	return data, server_http.Request(serverURL, ep, nil, &data, options.GetIdentity(), filesOp.logfile)
//}
//
//func (filesOp *filesHTTP) Remove(bucketID files.BucketID, path string, identity *auth.Identity) error {
//	ep := filesOp.endpoints[files.EPRemove]
//	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path
//
//	return server_http.Request(serverURL, ep, nil, nil, options.GetIdentity(), filesOp.logfile)
//}
//
//func (filesOp *filesHTTP) List(bucketID files.BucketID, path string, depth int, identity *auth.Identity) (files.FilesInfo, error) {
//	if path = strings.TrimSpace(path); path == "" {
//		path = "."
//	}
//
//	ep := filesOp.endpoints[files.EPList]
//	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + strconv.Itoa(depth)
//
//	var filesInfo files.FilesInfo
//	return filesInfo, server_http.Request(serverURL, ep, nil, &filesInfo, options.GetIdentity(), filesOp.logfile)
//}
//
//func (filesOp *filesHTTP) Stat(bucketID files.BucketID, path string, depth int, identity *auth.Identity) (*files.FileInfo, error) {
//	if path = strings.TrimSpace(path); path == "" {
//		path = "."
//	}
//
//	ep := filesOp.endpoints[files.EPStat]
//	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + strconv.Itoa(depth)
//
//	var fileInfo *files.FileInfo
//	return fileInfo, server_http.Request(serverURL, ep, nil, &fileInfo, options.GetIdentity(), filesOp.logfile)
//}
