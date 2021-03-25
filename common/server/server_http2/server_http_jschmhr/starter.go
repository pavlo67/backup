package server_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/starter"
	server_http "github.com/pavlo67/tools/common/server/server_http2"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	config       server.Config
	htmlTemplate string

	interfaceKey joiner.InterfaceKey
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Prepare(cfg *config.Config, options common.Map) error {
	ss.htmlTemplate = options.StringDefault("html_template", "")
	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	configKey := options.StringDefault("config_key", "server_http")
	if err := cfg.Value(configKey, &ss.config); err != nil {
		return err
	}

	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.OperatorV2 with key %s", logger.InterfaceKey)
	}

	onRequest, _ := joinerOp.Interface(server_http.OnRequestMiddlewareInterfaceKey).(server_http.OnRequestMiddleware)
	//if onRequest == nil {
	//	return fmt.Errorf("no server_http.OnRequestMiddleware with key %s", server_http.OnRequestMiddlewareInterfaceKey)
	//}

	wrappersHTTP := map[server_http.WrapperHTTPKey]WrapperHTTP{
		server_http.WrapperHTTPREST: WrapperHTTPREST,
	}
	if ss.htmlTemplate != "" {
		wrappersHTTP[server_http.WrapperHTTPPage] = WrapperHTTPPage(ss.htmlTemplate)
	}

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, onRequest, wrappersHTTP)
	if err != nil {
		return errors.Wrap(err, "on server_http_jschmhr.New()")
	}

	if err = joinerOp.Join(srvOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *serverHTTPJschmhr{} as server_http.OperatorV2 with key '%s'", ss.interfaceKey)
	}

	return nil

}
