package records_html

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

const InterfaceKey joiner.InterfaceKey = "records_html"

func Starter() starter.Operator {
	return &recordsHTMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsHTMLStarter{}

type recordsHTMLStarter struct {
	pagesConfig server_http.Config

	restConfig server_http.Config

	interfaceKey joiner.InterfaceKey
}

func (frhs *recordsHTMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (frhs *recordsHTMLStarter) Prepare(cfg *config.Config, options common.Map) error {
	var ok bool

	if frhs.pagesConfig, ok = options["pages_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["pages_config"]`)
	}
	if frhs.restConfig, ok = options["rest_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["rest_config"]`)
	}

	frhs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (frhs *recordsHTMLStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	formatterRecordsOp, err := New(frhs.pagesConfig, frhs.restConfig)
	if err != nil {
		return errors.CommonError(err, "can't init *formatterRecordsHTML as formatter.Operator")
	}

	if err = joinerOp.Join(formatterRecordsOp, frhs.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *recordsHTML as records_html.Operator with key '%s'", frhs.interfaceKey))
	}

	return nil
}
