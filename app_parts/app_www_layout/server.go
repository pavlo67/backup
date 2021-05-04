package app_www_layout

import (
	"io/ioutil"
	"os"

	"github.com/pavlo67/tools/common/server/server_http_v2/server_http_v2_jschmhr/wrapper_page"

	"github.com/pavlo67/tools/common/actor_www"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/common/server/server_http_v2/server_http_v2_jschmhr"
	"github.com/pavlo67/tools/common/thread"
)

func Serve(cfgService config.Config, l logger.Operator) (server_http_v2.OperatorV2, thread.KV) {
	commonChannel, err := thread.NewKV(&MenuWWW{})
	if err != nil {
		l.Fatalf("on thread.NewKV(): %s", err)
	}

	var commonFragments wrapper_page.CommonFragments = &SetMenu{Process: commonChannel}

	templatePath := filelib.CurrentPath() + "../../app_parts/app_www_layout/templates/local.html"
	htmlTemplateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		l.Fatalf("reading template (%s): %s", templatePath, err)
	}

	staticPath := filelib.CurrentPath() + "../../app_parts/app_www_layout/static/"
	fileInfo, err := os.Stat(staticPath)
	if err != nil {
		l.Fatalf("checking static files path (%s): %s", staticPath, err)
	}
	if !fileInfo.IsDir() {
		l.Fatalf("static files path (%s) is not a directory", staticPath)
	}

	//appMenu, _ := joinerOp.Interface(actor_www.AppMenuInterfaceKey).(thread.KVGetString)
	//, "app_menu": commonChannel

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{server_http_v2_jschmhr.Starter(), common.Map{"html_template": string(htmlTemplateBytes), string(actor_www.CommonFragmentsInterfaceKey): commonFragments}},
	}

	joinerOp, err := starter.Run(starters, &cfgService, "NB/WWW BUILD", l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http_v2.OperatorV2)
	if srvOp == nil {
		l.Fatalf("no server_http_v2.OperatorV2 with key %s", server_http.InterfaceKey)
	}

	if err = srvOp.HandleFiles("static", "/static/", server_http.StaticPath{LocalPath: staticPath}); err != nil {
		l.Fatalf(`on srvOp.HandleFiles("static", "/static/", server_http.StaticPath{LocalPath: %s}): got %s`, staticPath, err)
	}

	return srvOp, commonChannel

}
