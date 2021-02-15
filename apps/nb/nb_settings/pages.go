package nb_settings

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
)

// Swagger-UI sorts interface sections due to the first their path occurrences, so:
// 1. unauthorized   /auth/...
// 2. admin          /front/...

var restConfig = server_http.Config{
	Title:   "Notebook REST API",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}},

		//notebook.IntefaceKeyRESTRead:     {Path: "/read", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTChildren: {Path: "/children", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTTags:     {Path: "/tags", Tags: []string{"unauthorized"}},
		//notebook.IntefaceKeyRESTTagged:   {Path: "/tagged", Tags: []string{"unauthorized"}},
		//
		//notebook.IntefaceKeyRESTSave:   {Path: "/save", Tags: []string{"authorized"}},
		//notebook.IntefaceKeyRESTDelete: {Path: "/delete", Tags: []string{"authorized"}},
	},
}

var pagesConfig = server_http.Config{
	Title:   "Notebook pages",
	Version: "0.0.1",
	EndpointsSettled: map[joiner.InterfaceKey]server_http.EndpointSettled{
		notebook.IntefaceKeyHTMLRoot:   {Path: "", Tags: []string{"unauthorized"}},
		notebook.IntefaceKeyHTMLView:   {Path: "/view", Tags: []string{"unauthorized"}},
		notebook.IntefaceKeyHTMLTags:   {Path: "/tags", Tags: []string{"unauthorized"}},
		notebook.IntefaceKeyHTMLTagged: {Path: "/tagged", Tags: []string{"unauthorized"}},

		notebook.IntefaceKeyHTMLEdit: {Path: "/edit", Tags: []string{"unauthorized"}},
	},
}
