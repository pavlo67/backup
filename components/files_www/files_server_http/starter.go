package files_server_http

import (
	"fmt"

	"github.com/pavlo67/tools/components/files_www/files_server_http/files_html"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/entities/files"
)

const InterfaceKey joiner.InterfaceKey = "files_server_http"

func Starter() starter.Operator {
	return &filesServerHTTPStarter{}
}

var _ starter.Operator = &filesServerHTTPStarter{}

var l logger.Operator
var filesOp files.Operator
var filesHTMLOp files_html.Operator

type filesServerHTTPStarter struct {
	filesKey     joiner.InterfaceKey
	filesHTMLKey joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
}

func (fshs *filesServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fshs *filesServerHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {
	fshs.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	fshs.filesHTMLKey = joiner.InterfaceKey(options.StringDefault("files_html_key", string(files_html.InterfaceKey)))
	fshs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

const onRun = "on filesServerHTTPStarter.Execute()"

func (fshs *filesServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if filesOp, _ = joinerOp.Interface(fshs.filesKey).(files.Operator); filesOp == nil {
		return fmt.Errorf(onRun+": no files.Operator with key %s", fshs.filesKey)
	}

	if filesHTMLOp, _ = joinerOp.Interface(fshs.filesHTMLKey).(files_html.Operator); filesHTMLOp == nil {
		return fmt.Errorf("no files_html.Operator with key %s", fshs.filesHTMLKey)
	}

	return nil
}
