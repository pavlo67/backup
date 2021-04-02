package catalogue_server_http

import (
	"fmt"

	"github.com/pavlo67/tools/common/actor"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/entities/items"
)

const InterfaceKey joiner.InterfaceKey = "catalogue_server_http"

func Starter() starter.Operator {
	return &catalogueServerHTTPStarter{}
}

var _ starter.Operator = &catalogueServerHTTPStarter{}

var l logger.Operator

type catalogueServerHTTPStarter struct {
	prefix string

	itemsKey     joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
}

func (fshs *catalogueServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fshs *catalogueServerHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {
	fshs.prefix = options.StringDefault("prefix", "")

	fshs.itemsKey = joiner.InterfaceKey(options.StringDefault("catalogue_key", string(items.InterfaceKey)))
	fshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

const onRun = "on catalogueServerHTTPStarter.Execute()"

func (fshs *catalogueServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	itemsOp, _ := joinerOp.Interface(fshs.itemsKey).(items.Operator)
	if itemsOp == nil {
		return fmt.Errorf(onRun+": no items.Operator with key %s", fshs.itemsKey)
	}

	cp, err := newCataloguePages(fshs.prefix, itemsOp)
	if cp == nil || err != nil {
		return fmt.Errorf(onRun+": can't newCataloguePages(), got %#v / %s", cp, err)
	}

	if err := joinerOp.Join(cp, actor.ConfigPages); err != nil {
		return fmt.Errorf(onRun+": can't join *configPages with key %s, got %s", actor.ConfigPages, err)
	}

	return nil
}
