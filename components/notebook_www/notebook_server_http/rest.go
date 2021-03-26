package notebook_server_http

//var restPrefix = "/rest"
//
//// Swagger-UI sorts interface sections due to the first their path occurrences, so:
//// 1. unauthorized   /auth/...
//// 2. admin          /front/...
//
//var RestConfig = server_http.ConfigREST{
//
//	Title:     "Notebook REST API",
//	Version:   "0.0.1",
//	EndpointsSettled: server_http.EndpointsConfig{
//		// auth.IntefaceKeyAuthenticate:     {Path: "/auth", Tags: []string{"unauthorized"}},
//
//		// notebook.IntefaceKeyRESTRead:     {Path: "/read", Tags: []string{"unauthorized"}},
//		// notebook.IntefaceKeyRESTChildren: {Path: "/children", Tags: []string{"unauthorized"}},
//		// notebook.IntefaceKeyRESTTags:     {Path: "/tags", Tags: []string{"unauthorized"}},
//		// notebook.IntefaceKeyRESTList:     {Path: "/tagged", Tags: []string{"unauthorized"}},
//		//
//		// notebook.IntefaceKeyRESTSave:     {Path: "/save", Tags: []string{"authorized"}},
//		// notebook.IntefaceKeyRESTDele:     {Path: "/delete", Tags: []string{"authorized"}},
//	},
//}
//
//func HandleREST(joinerOp joiner.Operator, srvOp server_http.OperatorV2) error {
//
//	srvPort, isHTTPS := srvOp.Addr() // isHTTPS
//
//	// REST -----------------------------------------------------------
//
//	if err := RestConfig.Complete("", srvPort, restPrefix); err != nil {
//		return err
//	}
//	if err := RestConfig.HandlePages(srvOp, l); err != nil {
//		return err
//	}
//
//	// swagger.json ---------------------------------------------------
//
//	restStaticPath := filelib.CurrentPath() + "../rest_static/"
//	if err := RestConfig.InitSwagger(isHTTPS, restStaticPath+"api-docs/swagger.json", l); err != nil {
//		return err
//	}
//	//if err := srvOp.HandleFiles("swagger.json", restPrefix+"/swagger.json", server_http.StaticPath{LocalPath: restStaticPath + "api-docs/swagger.json"}); err != nil {
//	//	return err
//	//}
//
//	return nil
//}
//
//func HandleSwagger(joinerOp joiner.Operator, srvOp server_http.OperatorV2) error {
//
//	restStaticPath := filelib.CurrentPath() + "../rest_static/"
//	if err := srvOp.HandleFiles("rest_static", restPrefix+"/*filepath", server_http.StaticPath{LocalPath: restStaticPath}); err != nil {
//		return err
//	}
//
//	return nil
//}
