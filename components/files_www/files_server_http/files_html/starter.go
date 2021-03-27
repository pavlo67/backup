package files_html

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

const InterfaceKey joiner.InterfaceKey = "files_html"

func Starter() starter.Operator {
	return &filesHTMLStarter{}
}

var l logger.Operator
var _ starter.Operator = &filesHTMLStarter{}

type filesHTMLStarter struct {
	pagesConfig *server_http.ConfigPages
	// restConfig  *server_http.Config

	interfaceKey joiner.InterfaceKey
}

func (nhs *filesHTMLStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nhs *filesHTMLStarter) Prepare(cfg *config.Config, options common.Map) error {
	switch v := options["pages_config"].(type) {
	case server_http.ConfigPages:
		nhs.pagesConfig = &v
	case *server_http.ConfigPages:
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

func (nhs *filesHTMLStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.OperatorV2 with key %s", logger.InterfaceKey)
	}

	filesOp, err := New(*nhs.pagesConfig) // *nhs.restConfig
	if err != nil {
		return errors.CommonError(err, "can't init *filesHTML as files_html.Operator")
	}

	if err = joinerOp.Join(filesOp, nhs.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *filesHTML as files_html.Operator with key '%s'", nhs.interfaceKey))
	}

	return nil
}
