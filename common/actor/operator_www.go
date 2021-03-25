package actor

import (
	"sync"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	server_http "github.com/pavlo67/tools/common/server/server_http2"
	"github.com/pavlo67/tools/common/server/server_http2/server_http_jschmhr"
)

type OperatorWWW interface {
	Name() string
	Starters(options common.Map) ([]starter.Starter, error)
	Config() (server_http.Config, error)
	//Root() *server_http.Endpoint
	//Details() *server_http.Endpoint
	//Accept() *server_http.Endpoint
	//Search() *server_http.Endpoint
}

func RunOneWWW(srvOp server_http.OperatorV2, actorWWW OperatorWWW, cfgService *config.Config, options common.Map, l logger.Operator) {
	starters, err := actorWWW.Starters(options)
	if err != nil {
		l.Fatal(err)
	}

	joinerOpActor, err := starter.Run(starters, cfgService, "ACTOR BUILD: "+actorWWW.Name(), l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOpActor.CloseAll()

	pagesPrefix := actorWWW.Name()

	serverConfig, err := actorWWW.Config()
	if err != nil {
		l.Fatal(err)
	}

	port, _ := srvOp.Addr()

	serverConfig.Complete("", port, pagesPrefix)
	if err := serverConfig.HandlePages(srvOp, l); err != nil {
		l.Fatal(err)
	}

	//pagesStaticPath := filelib.CurrentPath() + "../pages_static/"
	//if err := srvOp.HandleFiles("pages_static", pagesPrefix+"/static/*filepath", server_http.StaticPath{LocalPath: pagesStaticPath}); err != nil {
	//	l.Fatal(err)
	//}

}

func RunWWW(cfgService *config.Config, htmlTemplate, label string, actorsWWW []OperatorWWW, l logger.Operator) {
	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
	}

	joinerOp, err := starter.Run(starters, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.OperatorV2)
	if srvOp == nil {
		l.Fatalf("no server_http.OperatorV2 with key %s", server_http.InterfaceKey)
	}

	options := common.Map{
		"html_template": htmlTemplate,
	}

	for _, actorWWW := range actorsWWW {
		RunOneWWW(srvOp, actorWWW, cfgService, options, l)
	}

	//if err := HandleSwagger(joinerOp, srvOp); err != nil {
	//	return err
	//}

	var WG sync.WaitGroup
	WG.Add(1)

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	WG.Wait()
}
