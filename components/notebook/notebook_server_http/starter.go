package notebook_server_http

import (
	"fmt"

	"github.com/pavlo67/tools/components/formatter"
	"github.com/pavlo67/tools/components/records/formatter_records_html"

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
	filesKey            joiner.InterfaceKey
	recordsKey          joiner.InterfaceKey
	formatterRecordsKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var recordsOp records.Operator
var filesOp files.Operator
var formatterRecordsOp formatter.Operator

func (nshs *notebookServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nshs *notebookServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {
	nshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	nshs.recordsKey = joiner.InterfaceKey(options.StringDefault("records_key", string(records.InterfaceKey)))
	nshs.formatterRecordsKey = joiner.InterfaceKey(options.StringDefault("formatter_records_key", string(formatter_records_html.InterfaceKey)))

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

	if formatterRecordsOp, _ = joinerOp.Interface(nshs.formatterRecordsKey).(formatter.Operator); formatterRecordsOp == nil {
		return fmt.Errorf("no formatter.Operator with key %s", nshs.formatterRecordsKey)
	}

	return server_http.JoinEndpoints(joinerOp, Endpoints)
}
