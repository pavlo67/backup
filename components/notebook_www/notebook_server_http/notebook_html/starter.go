package notebook_html

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	server_http "github.com/pavlo67/tools/common/server/server_http2"
)

const InterfaceKey joiner.InterfaceKey = "notebook_html"

func Starter() starter.Operator {
	return &notebookHTMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &notebookHTMLStarter{}

type notebookHTMLStarter struct {
	htmlTemplate string
	pagesConfig  *server_http.Config
	restConfig   *server_http.Config

	interfaceKey joiner.InterfaceKey
}

func (nhs *notebookHTMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nhs *notebookHTMLStarter) Prepare(cfg *config.Config, options common.Map) error {
	var ok bool

	if nhs.htmlTemplate, ok = options.String("html_template"); !ok {
		return fmt.Errorf(`no html template in options["html_template"]`)
	}

	switch v := options["pages_config"].(type) {
	case server_http.Config:
		nhs.pagesConfig = &v
	case *server_http.Config:
		nhs.pagesConfig = v
	}
	if nhs.pagesConfig == nil {
		return fmt.Errorf(`no server_http.Config in options["pages_config"]`)
	}

	//switch v := options["rest_config"].(type) {
	//case server_http.Config:
	//	nhs.restConfig = &v
	//case *server_http.Config:
	//	nhs.restConfig = v
	//}
	//if nhs.restConfig == nil {
	//	return fmt.Errorf(`no server_http.Config in options["rest_config"]`)
	//}

	nhs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (nhs *notebookHTMLStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.OperatorV2 with key %s", logger.InterfaceKey)
	}

	notebookOp, err := New(nhs.htmlTemplate, *nhs.pagesConfig) // *nhs.restConfig
	if err != nil {
		return errors.CommonError(err, "can't init *notebookHTML as notebook_html.OperatorV2")
	}

	if err = joinerOp.Join(notebookOp, nhs.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *notebookHTML as notebook_html.OperatorV2 with key '%s'", nhs.interfaceKey))
	}

	return nil
}
