package notebook_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/common/files"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http/notebook_html"
	"github.com/pavlo67/tools/entities/records"
)

const InterfaceKey joiner.InterfaceKey = "notebook_server_http"

func Starter() starter.Operator {
	return &notebookServerHTTPStarter{}
}

var _ starter.Operator = &notebookServerHTTPStarter{}

type notebookServerHTTPStarter struct {
	filesKey    joiner.InterfaceKey
	recordsKey  joiner.InterfaceKey
	tagsHTMLKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var recordsOp records.Operator

// var filesOp files.Operator
var notebookHTMLOp *notebook_html.HTMLOp

func (nshs *notebookServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nshs *notebookServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {
	nshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	nshs.recordsKey = joiner.InterfaceKey(options.StringDefault("records_key", string(records.InterfaceKey)))

	nshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

const onRun = "on notebookServerHTTPStarter.Run()"

func (nshs *notebookServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf(onRun+": no logger.Operator with key %s", logger.InterfaceKey)
	}

	//if filesOp, _ = joinerOp.Interface(nshs.filesKey).(files.Operator); filesOp == nil {
	//	return fmt.Errorf(onRun+": no files.Operator with key %s", nshs.filesKey)
	//}

	if recordsOp, _ = joinerOp.Interface(nshs.recordsKey).(records.Operator); recordsOp == nil {
		return fmt.Errorf(onRun+": no records.Operator with key %s", nshs.recordsKey)
	}

	var err error
	if notebookHTMLOp, err = notebook_html.New(PagesConfig, l); err != nil {
		return fmt.Errorf(onRun+": %s", err)
	}

	return nil
}
