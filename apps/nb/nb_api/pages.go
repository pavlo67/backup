package nb_api

import (
	"errors"
	"fmt"

	"github.com/pavlo67/tools/components/notebook"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
)

// Swagger-UI sorts interface sections due to the first their path occurrences, so:
// 1. unauthorized   /auth/...
// 2. admin          /front/...

// TODO!!! keep in mind that EndpointsConfig key and corresponding .HandlerKey not necessarily are the same, they can be defined different

var restConfig = server_http.Config{
	Title:   "Notebook REST API",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}, EndpointInternalKey: auth.IntefaceKeyAuthenticate},
	},
}

var pagesConfig = server_http.Config{
	Title:   "Notebook pages",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		notebook.IntefaceKeyRoot: {Path: "", Tags: []string{"unauthorized"}, EndpointInternalKey: notebook.IntefaceKeyRoot},
	},
}

const onInitPages = "on nb_api.InitPages()"

func InitPages(srvOp server_http.Operator, pagesCfg server_http.Config, l logger.Operator) error {
	if srvOp == nil {
		return errors.New(onInitPages + ": srvOp == nil")
	}

	for key, ep := range pagesCfg.EndpointsSettled {
		if err := srvOp.HandleEndpoint(key, pagesCfg.Prefix+ep.Path, ep.Endpoint); err != nil {
			return fmt.Errorf(onInitPages+": handling %s, %s, %#v got %s", key, ep.Path, ep.Endpoint, err)
		}
	}

	return nil
	// return srvOp.HandleFiles("swagger", pagesCfg.Prefix+"/"+swaggerSubpath+"/*filepath", StaticPath{LocalPath: swaggerPath})
}
