package nb_settings

import (
	"fmt"
	"io/ioutil"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http"
)

// Swagger-UI sorts interface sections due to the first their path occurrences, so:
// 1. unauthorized   /auth/...
// 2. admin          /front/...

var RestConfig = server_http.Config{
	Title:   "Notebook REST API",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}},

		//notebook.IntefaceKeyRESTRead:     {Path: "/read", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTChildren: {Path: "/children", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTTags:     {Path: "/tags", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTList:   {Path: "/tagged", Tags: []string{"unauthorized"}},
		//
		//notebook.IntefaceKeyRESTSave:   {Path: "/save", Tags: []string{"authorized"}},
		//notebook.IntefaceKeyRESTDele: {Path: "/delete", Tags: []string{"authorized"}},
	},
}

var PagesConfig = server_http.Config{
	Title:   "Notebook pages",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		notebook.IntefaceKeyHTMLRoot: {Path: ""},
		notebook.IntefaceKeyHTMLView: {Path: "/view"},
		notebook.IntefaceKeyHTMLTags: {Path: "/tags"},
		notebook.IntefaceKeyHTMLList: {Path: "/list"},
		notebook.IntefaceKeyHTMLEdit: {Path: "/edit"},
		notebook.IntefaceKeyHTMLSave: {Path: "/save"},
	},
}

var pagesPrefix = ""
var restPrefix = "/rest"

func CompleteServerConfigs() (string, *server_http.Config, *server_http.Config, error) {
	templatePath := filelib.CurrentPath() + "../templates/local.html"
	htmlTemplateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("on ioutil.ReadFile(%s): %s", templatePath, err)
	}

	if err := PagesConfig.CompleteDirectly(notebook_server_http.Pages, "", 0, pagesPrefix); err != nil {
		return "", nil, nil, fmt.Errorf(`on PagesConfig.CompleteDirectly() got %s`, err)
	}

	if err := RestConfig.CompleteDirectly(auth_server_http.Endpoints, "", 0, restPrefix); err != nil {
		return "", nil, nil, fmt.Errorf(`on RestConfig.CompleteDirectly() got %s`, err)
	}

	return string(htmlTemplateBytes), &PagesConfig, &RestConfig, nil
}

func HandlePages(joinerOp joiner.Operator, srvOp server_http.Operator) error {

	srvPort, isHTTPS := srvOp.Addr() // isHTTPS

	if err := RestConfig.CompleteWithJoiner(joinerOp, "", srvPort, restPrefix); err != nil {
		return err
	}
	if err := server_http.InitPages(srvOp, RestConfig, l); err != nil {
		return err
	}
	if err := PagesConfig.CompleteWithJoiner(joinerOp, "", srvPort, pagesPrefix); err != nil {
		return err
	}
	if err := server_http.InitPages(srvOp, PagesConfig, l); err != nil {
		return err
	}

	restStaticPath := filelib.CurrentPath() + "../rest_static/"
	restServerPath := restPrefix + "/*filepath"

	swaggerStaticPath := restStaticPath + "api-docs/"
	swaggerStaticFilePath := swaggerStaticPath + "swaggerJSON.json"
	swaggerJSON, err := RestConfig.SwaggerV2(isHTTPS)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(swaggerStaticFilePath, swaggerJSON, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerStaticFilePath, swaggerJSON, err)
	}
	l.Infof("%d bytes are written into %s", len(swaggerJSON), swaggerStaticFilePath)

	if err := srvOp.HandleFiles("rest_static", restServerPath, server_http.StaticPath{LocalPath: restStaticPath}); err != nil {
		return err
	}

	pagesStaticPath := filelib.CurrentPath() + "../pages_static/"
	pagesStaticServerPath := pagesPrefix + "/static/*filepath"
	if pagesStaticPath != "" {
		if err := srvOp.HandleFiles("pages_static", pagesStaticServerPath, server_http.StaticPath{LocalPath: pagesStaticPath}); err != nil {
			return err
		}
	}

	return nil
}
