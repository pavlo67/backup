package nb_settings

import (
	"fmt"
	"sync"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &nbStarter{}
}

var _ starter.Operator = &nbStarter{}

type nbStarter struct{}

// --------------------------------------------------------------------------

var WG sync.WaitGroup
var l logger.Operator

func (ns *nbStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ns *nbStarter) Prepare(cfg *config.Config, options common.Map) error {
	var cfgStorage common.Map
	if err := cfg.Value("nb_api", &cfgStorage); err != nil {
		return errors.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	}

	return nil
}

func (ns *nbStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if srvOp == nil {
		return fmt.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	if err := HandlePages(joinerOp, srvOp); err != nil {
		return err
	}

	WG.Add(1)

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	return nil
}

// TODO!!! customize it
// if isHTTPS {
//	go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
// }
