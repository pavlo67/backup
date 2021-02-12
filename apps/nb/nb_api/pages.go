package nb_api

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/tools/components/notebook"

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
		auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}, Endpoint: server_http.Endpoint{InternalKey: auth.IntefaceKeyAuthenticate}},
	},
}

var pagesConfig = server_http.Config{
	Title:   "Notebook pages",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		notebook.IntefaceKeyRoot: {Path: "", Tags: []string{"unauthorized"}, Endpoint: server_http.Endpoint{InternalKey: notebook.IntefaceKeyRoot}},
		notebook.IntefaceKeyEdit: {Path: "/view", Tags: []string{"unauthorized"}, Endpoint: server_http.Endpoint{InternalKey: notebook.IntefaceKeyEdit}},
		notebook.IntefaceKeyView: {Path: "/edit", Tags: []string{"unauthorized"}, Endpoint: server_http.Endpoint{InternalKey: notebook.IntefaceKeyView}},
	},
}
