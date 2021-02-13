package files_fs

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/components/files"
	"github.com/pkg/errors"
)

func Starter() starter.Operator {
	return &filesFSStarter{}
}

var l logger.Operator
var _ starter.Operator = &filesFSStarter{}

type filesFSStarter struct {
	buckets      files.Buckets
	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey

	// pathInfix    string
}

func (ffs *filesFSStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ffs *filesFSStarter) Prepare(cfg *config.Config, options common.Map) error {

	ffs.buckets, _ = options["buckets"].(files.Buckets)
	if ffs.buckets == nil {
		return errors.Errorf("no 'buckets' in options: %#v", options)
	}

	//configKey := strings.TrimSpace(options.StringDefault("config_key", "buckets"))
	//if configKey == "" {
	//	return nil, errors.Errorf("no 'config_key' in options (%#v)", options)
	//}
	//if err := cfg.Value(configKey, &ffs.buckets); err != nil {
	//	l.Errorf("1111111111 in config: %#v", cfg)
	//
	//	return nil, errors.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	//}

	ffs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(files.InterfaceKey)))
	if ffs.interfaceKey == "" {
		return errors.Errorf("no 'interface_key' in options (%#v)", options)
	}

	ffs.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(files.InterfaceKeyCleaner)))
	if ffs.cleanerKey == "" {
		return errors.Errorf("no 'cleaner_key' in options (%#v)", options)
	}

	return nil
}

func (ffs *filesFSStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	filesOp, filesCleanerOp, err := New(ffs.buckets)
	if err != nil {
		return errors.Wrap(err, "can't init *filesFS{} as files.Operator")
	}

	if err = joinerOp.Join(filesOp, ffs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *filesFS{} as files.Operator with key '%s'", ffs.interfaceKey)
	}

	if err = joinerOp.Join(filesCleanerOp, ffs.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *filesFS{} as crud.Cleaner with key '%s'", ffs.cleanerKey)
	}

	return nil
}
