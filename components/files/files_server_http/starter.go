package files_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
)

const InterfaceKey joiner.InterfaceKey = "files_server_http"

func Starter() starter.Operator {
	return &managementServerHTTPStarter{}
}

var _ starter.Operator = &managementServerHTTPStarter{}

var l logger.Operator

var filesOp files.Operator

type managementServerHTTPStarter struct {
	filesKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

func (mshs *managementServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (mshs *managementServerHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {
	mshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	mshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

const onRun = "on managementServerHTTPStarter.Execute()"

func (mshs *managementServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if filesOp, _ = joinerOp.Interface(mshs.filesKey).(files.Operator); filesOp == nil {
		return fmt.Errorf(onRun+": no files.Operator with key %s", mshs.filesKey)
	}

	if err := joinerOp.Join(&readEndpoint, files.HandlerRead); err != nil {
		return errors.Wrapf(err, "can't join readEndpoint as server_http.Endpoint with key '%s'", files.HandlerRead)
	}
	if err := joinerOp.Join(&saveEndpoint, files.HandlerSave); err != nil {
		return errors.Wrapf(err, "can't join saveEndpoint as server_http.Endpoint with key '%s'", files.HandlerSave)
	}
	if err := joinerOp.Join(&removeEndpoint, files.HandlerRemove); err != nil {
		return errors.Wrapf(err, "can't join removeEndpoint as server_http.Endpoint with key '%s'", files.HandlerRemove)
	}
	if err := joinerOp.Join(&listEndpoint, files.HandlerList); err != nil {
		return errors.Wrapf(err, "can't join listEndpoint as server_http.Endpoint with key '%s'", files.HandlerList)
	}
	if err := joinerOp.Join(&statEndpoint, files.HandlerStat); err != nil {
		return errors.Wrapf(err, "can't join statEndpoint as server_http.Endpoint with key '%s'", files.HandlerStat)
	}

	return nil
}
