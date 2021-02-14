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
	"github.com/pavlo67/tools/components/records/records_html"
	"github.com/pavlo67/tools/components/tags/tags_html"
)

const InterfaceKey joiner.InterfaceKey = "notebook_server_http"

func Starter() starter.Operator {
	return &notebookServerHTTPStarter{}
}

var _ starter.Operator = &notebookServerHTTPStarter{}

type notebookServerHTTPStarter struct {
	filesKey       joiner.InterfaceKey
	recordsKey     joiner.InterfaceKey
	recordsHTMLKey joiner.InterfaceKey
	tagsHTMLKey    joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var recordsOp records.Operator
var filesOp files.Operator
var recordsHTMLOp records_html.Operator
var tagsHTMLOp tags_html.Operator

func (nshs *notebookServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nshs *notebookServerHTTPStarter) Prepare(_ *config.Config, options common.Map) error {
	nshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	nshs.recordsKey = joiner.InterfaceKey(options.StringDefault("records_key", string(records.InterfaceKey)))
	nshs.recordsHTMLKey = joiner.InterfaceKey(options.StringDefault("records_html_key", string(records_html.InterfaceKey)))
	nshs.tagsHTMLKey = joiner.InterfaceKey(options.StringDefault("tags_html_key", string(tags_html.InterfaceKey)))

	nshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (nshs *notebookServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if filesOp, _ = joinerOp.Interface(nshs.filesKey).(files.Operator); filesOp == nil {
		return fmt.Errorf("no files.Operator with key %s", nshs.filesKey)
	}

	if recordsOp, _ = joinerOp.Interface(nshs.recordsKey).(records.Operator); recordsOp == nil {
		return fmt.Errorf("no records.Operator with key %s", nshs.recordsKey)
	}

	if recordsHTMLOp, _ = joinerOp.Interface(nshs.recordsHTMLKey).(records_html.Operator); recordsHTMLOp == nil {
		return fmt.Errorf("no records_html.Operator with key %s", nshs.recordsHTMLKey)
	}
	if tagsHTMLOp, _ = joinerOp.Interface(nshs.tagsHTMLKey).(tags_html.Operator); tagsHTMLOp == nil {
		return fmt.Errorf("no tags_html.Operator with key %s", nshs.tagsHTMLKey)
	}

	return server_http.JoinEndpoints(joinerOp, Endpoints)
}
