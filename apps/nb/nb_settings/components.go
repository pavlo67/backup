package nb_settings

import (
	"fmt"

	"github.com/pavlo67/common/common/auth/auth_http"
	"github.com/pavlo67/tools/components/files/files_http"
	"github.com/pavlo67/tools/components/records/records_http"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_persons"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/persons/persons_fs"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/connect/connect_sqlite"
	"github.com/pavlo67/tools/components/files"
	"github.com/pavlo67/tools/components/files/files_fs"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http"
	"github.com/pavlo67/tools/components/records/records_html"
	"github.com/pavlo67/tools/components/records/records_sqlite"
	"github.com/pavlo67/tools/components/tags/tags_html"
)

// TODO!!!
var bucketsOptions = common.Map{
	"buckets": files.Buckets{files.BucketID("1"): "../1"},
}

var pagesPrefix = ""
var restPrefix = "/rest"

var prefixOptions = common.Map{
	"pages_prefix": pagesPrefix,
	"rest_prefix":  restPrefix,
}

func ServerComponents() ([]starter.Starter, error) {

	if err := pagesConfig.CompleteDirectly(notebook_server_http.Pages, "", 0, pagesPrefix); err != nil {
		return nil, fmt.Errorf(`on pagesConfig.CompleteDirectly() got %s`, err)
	}
	if err := restConfig.CompleteDirectly(auth_server_http.Endpoints, "", 0, restPrefix); err != nil {
		return nil, fmt.Errorf(`on restConfig.CompleteDirectly() got %s`, err)
	}

	endpointsOptions := common.Map{
		"pages_config": pagesConfig,
		"rest_config":  restConfig,
	}

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{connect_sqlite.Starter(), nil},

		// auth/persons components
		{persons_fs.Starter(), nil},
		{auth_persons.Starter(), nil},
		{auth_jwt.Starter(), nil},
		{auth_server_http.Starter(), nil}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		{files_fs.Starter(), bucketsOptions},
		{records_sqlite.Starter(), nil},
		{records_html.Starter(), endpointsOptions},
		{tags_html.Starter(), endpointsOptions},
		{notebook_server_http.Starter(), nil},

		// action managers
		{server_http_jschmhr.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		{Starter(), prefixOptions},
	}

	return starters, nil
}

func ClientComponents() ([]starter.Starter, error) {

	if err := pagesConfig.CompleteDirectly(notebook_server_http.Pages, "", 0, pagesPrefix); err != nil {
		return nil, fmt.Errorf(`on pagesConfig.CompleteDirectly() got %s`, err)
	}
	if err := restConfig.CompleteDirectly(auth_server_http.Endpoints, "", 0, restPrefix); err != nil {
		return nil, fmt.Errorf(`on restConfig.CompleteDirectly() got %s`, err)
	}

	endpointsOptions := common.Map{
		"pages_config": pagesConfig,
		"rest_config":  restConfig,
	}

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},

		// auth/persons components
		{auth_jwt.Starter(), nil},
		{auth_http.Starter(), nil}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		{files_http.Starter(), endpointsOptions},
		{records_http.Starter(), endpointsOptions},
	}

	return starters, nil
}
