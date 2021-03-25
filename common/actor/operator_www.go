package actor

import (
	"fmt"
	"sync"

	"github.com/pavlo67/common/common/joiner"

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
	Config() (server_http.ConfigPages, error)
	//Root() *server_http.EndpointREST
	//Details() *server_http.EndpointREST
	//Accept() *server_http.EndpointREST
	//Search() *server_http.EndpointREST
}

func RunOneWWW(srvOp server_http.OperatorV2, actorWWW OperatorWWW, cfgService *config.Config, options common.Map, l logger.Operator) (joiner.Operator, error) {
	starters, err := actorWWW.Starters(options)
	if err != nil {
		l.Fatal(err)
	}

	joinerOp, err := starter.Run(starters, cfgService, "ACTOR BUILD: "+actorWWW.Name(), l)
	if err != nil {
		return joinerOp, err
	}

	pagesPrefix := actorWWW.Name()

	serverConfig, err := actorWWW.Config()
	if err != nil {
		return joinerOp, err
	}

	port, _ := srvOp.Addr()

	serverConfig.Complete("", port, pagesPrefix)
	if err := serverConfig.HandlePages(srvOp, l); err != nil {
		return joinerOp, err
	}

	//pagesStaticPath := filelib.CurrentPath() + "../pages_static/"
	//if err := srvOp.HandleFiles("pages_static", pagesPrefix+"/static/*filepath", server_http.StaticPath{LocalPath: pagesStaticPath}); err != nil {
	//	l.Fatal(err)
	//}

	return joinerOp, nil

}

func RunWWW(cfgService *config.Config, htmlTemplate, label string, actorsWWW []OperatorWWW, l logger.Operator) ([]joiner.Operator, error) {

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"html_template": htmlTemplate}},
	}

	var joinerOps []joiner.Operator

	joinerOp, err := starter.Run(starters, cfgService, label, l)
	joinerOps = append(joinerOps, joinerOp)
	if err != nil {
		return joinerOps, err
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.OperatorV2)
	if srvOp == nil {
		return joinerOps, fmt.Errorf("no server_http.OperatorV2 with key %s", server_http.InterfaceKey)
	}

	for _, actorWWW := range actorsWWW {
		joinerOp, err := RunOneWWW(srvOp, actorWWW, cfgService, nil, l)
		joinerOps = append(joinerOps, joinerOp)
		if err != nil {
			return joinerOps, err
		}
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

	return joinerOps, nil
}
