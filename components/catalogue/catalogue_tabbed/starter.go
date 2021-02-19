package catalogue_tabbed

import (
	"fmt"

	"github.com/pavlo67/tools/components/files"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/catalogue"
)

func Starter() starter.Operator {
	return &catalogueTabbedStarter{}
}

var l logger.Operator
var _ starter.Operator = &catalogueTabbedStarter{}

type catalogueTabbedStarter struct {
	filesKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey

	// pathInfix    string
}

func (cts *catalogueTabbedStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cts *catalogueTabbedStarter) Prepare(cfg *config.Config, options common.Map) error {

	cts.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	cts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(catalogue.InterfaceKey)))
	cts.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(catalogue.InterfaceKeyCleaner)))

	return nil
}

func (cts *catalogueTabbedStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	filesOp, _ := joinerOp.Interface(cts.filesKey).(files.Operator)
	if filesOp == nil {
		return fmt.Errorf("no files.Operator with key %s", cts.filesKey)
	}

	catalogueOp, catalogueCleanerOp, err := New(filesOp)
	if err != nil {
		return errors.CommonError(err, "can't init *catalogueTabbed{} as catalogue.Operator")
	}

	if err = joinerOp.Join(catalogueOp, cts.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *catalogueTabbed{} as catalogue.Operator with key '%s'", cts.interfaceKey))
	}

	if err = joinerOp.Join(catalogueCleanerOp, cts.cleanerKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *catalogueTabbed{} as crud.Cleaner with key '%s'", cts.cleanerKey))
	}

	return nil
}
