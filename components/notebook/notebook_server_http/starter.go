package notebook_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
	"github.com/pavlo67/tools/components/records"
)

const InterfaceKey joiner.InterfaceKey = "notebook_server_http"

func Starter() starter.Operator {
	return &notebookServerHTTPStarter{}
}

var _ starter.Operator = &notebookServerHTTPStarter{}

type notebookServerHTTPStarter struct {
	filesKey   joiner.InterfaceKey
	recordsKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var recordsOp records.Operator
var filesOp files.Operator

func (nshs *notebookServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nshs *notebookServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {
	nshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	nshs.recordsKey = joiner.InterfaceKey(options.StringDefault("records_key", string(records.InterfaceKey)))

	nshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (nshs *notebookServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	// endpoints --------------------------------------------------------

	if filesOp, _ = joinerOp.Interface(nshs.filesKey).(files.Operator); filesOp == nil {
		return fmt.Errorf("no files.Operator with key %s", nshs.filesKey)
	}

	if recordsOp, _ = joinerOp.Interface(nshs.recordsKey).(records.Operator); recordsOp == nil {
		return fmt.Errorf("no records.Operator with key %s", nshs.recordsKey)
	}

	l.Infof("111111111 %#v", Endpoints)

	return server_http.JoinEndpoints(joinerOp, Endpoints)
}
