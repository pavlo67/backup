package actor

import (
	"fmt"
	"sync"

	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/common/server/server_http_v2/server_http_v2_jschmhr"
)

type OperatorWWW interface {
	Name() string
	Starters() ([]starter.Starter, error)
	Config() (*Config, *server_http.Config, *server_http.ConfigPages, error)

	//Root() *server_http.Endpoint
	//Details() *server_http.Endpoint
	//Accept() *server_http.Endpoint
	//Search() *server_http.Endpoint
}

func RunOneWWW(srvOp server_http.OperatorV2, actorWWW OperatorWWW, cfgService *config.Config, l logger.Operator) (joiner.Operator, error) {
	starters, err := actorWWW.Starters()
	if err != nil {
		l.Fatal(err)
	}

	joinerOp, err := starter.Run(starters, cfgService, "ACTOR BUILD: "+actorWWW.Name(), l)
	if err != nil {
		return joinerOp, err
	}

	actorConfig, endpointsConfig, pagesConfig, err := actorWWW.Config()
	if err != nil || actorConfig == nil {
		return joinerOp, fmt.Errorf("on actorWWW.Config(): got %#v / %#v / %#v / %s", actorConfig, endpointsConfig, pagesConfig, err)
	}

	port, _ := srvOp.Addr()
	pagesConfig.Complete("", port, actorConfig.Prefix)
	if err := pagesConfig.HandlePages(srvOp, l); err != nil {
		return joinerOp, err
	}

	return joinerOp, nil

}

func RunWWW(cfgService *config.Config, label, htmlTemplate, staticPath string, processMenu thread.FIFOKVItemsGetString, actorsWWW []OperatorWWW,
	l logger.Operator) (joinerOps []joiner.Operator, err error) {

	// initiating common components & http server -----------------------------------------

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{server_http_v2_jschmhr.Starter(), common.Map{"html_template": htmlTemplate, "process_menu": processMenu}},
	}

	//if err := HandleSwagger(joinerOp, srvOp); err != nil {
	//	return err
	//}

	joinerOp, err := starter.Run(starters, cfgService, label, l)
	joinerOps = append(joinerOps, joinerOp)
	if err != nil {
		return joinerOps, err
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.OperatorV2)
	if srvOp == nil {
		return joinerOps, fmt.Errorf("no server_http.OperatorV2 with key %s", server_http.InterfaceKey)
	}

	if err = srvOp.HandleFiles("static", "/static/", server_http.StaticPath{LocalPath: staticPath}); err != nil {
		return joinerOps, err
	}

	// initiating actors on http server ---------------------------------------------------

	for _, actorWWW := range actorsWWW {
		joinerOp, err := RunOneWWW(srvOp, actorWWW, cfgService, l)
		joinerOps = append(joinerOps, joinerOp)
		if err != nil {
			return joinerOps, err
		}
	}

	// starting http server ---------------------------------------------------------------

	var WG sync.WaitGroup
	WG.Add(1)

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	WG.Wait()

	// ------------------------------------------------------------------------------------

	return joinerOps, nil
}
