package storage_settings

import (
	"fmt"

	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &storageStarter{}
}

var l logger.Operator

var _ starter.Operator = &storageStarter{}

type storageStarter struct {
	prefix string
	// baseDir string

	// skipAbsentEndpoints bool
}

func (ss *storageStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *storageStarter) Prepare(cfg *config.Config, options common.Map) error {
	var cfgStorage common.Map
	if err := cfg.Value("storage_api", &cfgStorage); err != nil {
		return errors.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	}

	ss.prefix = cfgStorage.StringDefault("prefix", "")

	return nil
}

var serverConfig = server_http.Config{
	Title:            "Storage REST API",
	Version:          "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		// auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}, EndpointInternalKey: auth.IntefaceKeyAuthenticate},
	},
}

func (ss *storageStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if srvOp == nil {
		return fmt.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvPort, isHTTPS := srvOp.Addr()

	if err := serverConfig.CompleteWithJoiner(joinerOp, "", srvPort, ss.prefix); err != nil {
		return err
	}

	if err := server_http.InitEndpointsWithSwaggerV2(
		srvOp, serverConfig, !isHTTPS,
		filelib.CurrentPath()+"api-docs/", "api-docs",
		l,
	); err != nil {
		return err
	}

	WG.Add(1)

	// TODO!!! customize it
	// if isHTTPS {
	//	go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
	// }

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	return nil
}
