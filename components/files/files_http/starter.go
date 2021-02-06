package files_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
)

func Starter() starter.Operator {
	return &filesHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &filesHTTPStarter{}

type filesHTTPStarter struct {
	config config.Access

	prefix    string
	endpoints server_http.Endpoints

	interfaceKey joiner.InterfaceKey

	mockHandlers bool
	logfile      string
}

func (ihs *filesHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ihs *filesHTTPStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgHTTP config.Access
	if err := cfg.Value("files_http", &cfgHTTP); err != nil {
		return nil, err
	}

	ihs.config = cfgHTTP

	// TODO!!! pass for each server separately
	ihs.prefix = options.StringDefault("prefix", "")

	ihs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(files.InterfaceKey)))
	ihs.mockHandlers = options.IsTrue("mock_handlers")
	ihs.logfile = options.StringDefault("log_file", "")

	if endpoints, ok := options["endpoints"].(server_http.Endpoints); ok {
		ihs.endpoints = endpoints
	} else if endpointsPtr, ok := options["endpoints"].(*server_http.Endpoints); ok {
		ihs.endpoints = *endpointsPtr
	} else {
		return nil, errors.New("no endpoints description for filesHTTPStarter")
	}

	return nil, nil
}

func (ihs *filesHTTPStarter) Run(joinerOp joiner.Operator) error {

	filesOp, err := New(ihs.config, ihs.prefix, ihs.endpoints, ihs.mockHandlers, ihs.logfile)
	if err != nil {
		return errors.Wrap(err, "can't init *filesHTTP{} as files.Operator")
	}

	if err = joinerOp.Join(filesOp, ihs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *filesHTTP{} as files.Operator with key '%s'", ihs.interfaceKey)
	}

	return nil
}
