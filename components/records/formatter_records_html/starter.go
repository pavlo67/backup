package formatter_records_html

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "formatter_records"

func Starter() starter.Operator {
	return &formatterRecordsHTMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &formatterRecordsHTMLStarter{}

type formatterRecordsHTMLStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (frhs *formatterRecordsHTMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (frhs *formatterRecordsHTMLStarter) Prepare(cfg *config.Config, options common.Map) error {
	frhs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (frhs *formatterRecordsHTMLStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	formatterRecordsOp, err := New()
	if err != nil {
		return errata.CommonError(err, "can't init *formatterRecordsHTML as formatter.Operator")
	}

	if err = joinerOp.Join(formatterRecordsOp, frhs.interfaceKey); err != nil {
		return errata.CommonError(err, fmt.Sprintf("can't join *formatterRecordsHTML as formatter.Operator with key '%s'", frhs.interfaceKey))
	}

	return nil
}
