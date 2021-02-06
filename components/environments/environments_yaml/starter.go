package environments_yaml

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/environments"
)

func Starter() starter.Operator {
	return &environmentsYAMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &environmentsYAMLStarter{}

type environmentsYAMLStarter struct {
	config config.Access

	interfaceKey joiner.InterfaceKey
}

func (eys *environmentsYAMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (eys *environmentsYAMLStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	eys.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(environments.InterfaceKey)))
	if eys.interfaceKey == "" {
		return nil, errors.Errorf("no 'interface_key' in options (%#v)", options)
	}

	configKey := strings.TrimSpace(options.StringDefault("config_key", "environments_yaml"))
	if err := cfg.Value(configKey, &eys.config); err != nil {
		return nil, errata.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	}

	return nil, nil
}

func (eys *environmentsYAMLStarter) Run(joinerOp joiner.Operator) error {

	environmentsOp, err := New(eys.config)
	if err != nil {
		return errors.Wrap(err, "can't init *environmentsYAML{} as environments.Operator")
	}

	if err = joinerOp.Join(environmentsOp, eys.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *environmentsYAML{} as environments.Operator with key '%s'", eys.interfaceKey)
	}

	return nil
}
