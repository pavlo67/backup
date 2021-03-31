package catalogue_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/entities/catalogue"
)

const InterfaceKey joiner.InterfaceKey = "catalogue_server_http"

func Starter() starter.Operator {
	return &catalogueServerHTTPStarter{}
}

var _ starter.Operator = &catalogueServerHTTPStarter{}

var l logger.Operator
var catalogueOp catalogue.Operator
var filesHTMLOp *catalogueHTML

type catalogueServerHTTPStarter struct {
	prefix       string
	catalogueKey joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
}

func (fshs *catalogueServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fshs *catalogueServerHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {
	fshs.catalogueKey = joiner.InterfaceKey(options.StringDefault("catalogue_key", string(catalogue.InterfaceKey)))
	fshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	fshs.prefix = options.StringDefault("prefix", "")

	return nil
}

const onRun = "on catalogueServerHTTPStarter.Execute()"

func (fshs *catalogueServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if catalogueOp, _ = joinerOp.Interface(fshs.catalogueKey).(catalogue.Operator); catalogueOp == nil {
		return fmt.Errorf(onRun+": no catalogue.Operator with key %s", fshs.catalogueKey)
	}

	l.Infof("!!!!!!!!!! %s", fshs.prefix)

	var err error

	if filesHTMLOp, err = New(PagesConfig, fshs.prefix); err != nil {
		return fmt.Errorf(onRun+": %s", err)
	}

	return nil
}
