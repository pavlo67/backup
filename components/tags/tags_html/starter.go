package tags_html

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "tags_html"

func Starter() starter.Operator {
	return &tagsHTMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &tagsHTMLStarter{}

type tagsHTMLStarter struct {
	pagesConfig server_http.Config
	restConfig  server_http.Config

	interfaceKey joiner.InterfaceKey
}

func (fths *tagsHTMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fths *tagsHTMLStarter) Prepare(cfg *config.Config, options common.Map) error {
	var ok bool

	if fths.pagesConfig, ok = options["pages_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["pages_config"]`)
	}
	if fths.restConfig, ok = options["rest_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["rest_config"]`)
	}

	fths.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (fths *tagsHTMLStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	formatterRecordsOp, err := New(fths.pagesConfig, fths.restConfig)
	if err != nil {
		return errors.CommonError(err, "can't init *formatterRecordsHTML as formatter.Operator")
	}

	if err = joinerOp.Join(formatterRecordsOp, fths.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *tagsHTML as tags_html.Operator with key '%s'", fths.interfaceKey))
	}

	return nil
}
