package nb_www_settings

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth/auth_http"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/notebook/notebook_server_http"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http/notebook_html"
	"github.com/pavlo67/tools/entities/files"
	"github.com/pavlo67/tools/entities/files/files_fs"
	"github.com/pavlo67/tools/entities/files/files_http"
	"github.com/pavlo67/tools/entities/persons/auth_persons"
	"github.com/pavlo67/tools/entities/persons/persons_fs"
	"github.com/pavlo67/tools/entities/records/records_http"
	"github.com/pavlo67/tools/entities/records/records_sqlite"
)

// TODO!!!
var bucketsOptions = common.Map{
	"buckets": files.Buckets{files.BucketID("1"): "../1"},
}

func ServerComponents() ([]starter.Starter, error) {
	htmlTemplate, pagesConfig, restConfig, err := CompleteServerConfigs()
	if err != nil {
		return nil, err
	}

	renderOptions := common.Map{
		"html_template": htmlTemplate,
		"pages_config":  pagesConfig,
		"rest_config":   restConfig,
	}

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{db_sqlite.Starter(), nil},

		// auth/persons components
		{persons_fs.Starter(), nil},
		{auth_persons.Starter(), nil},
		{auth_jwt.Starter(), nil},
		{auth_server_http.Starter(), nil}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		{files_fs.Starter(), bucketsOptions},
		{records_sqlite.Starter(), nil},
		{notebook_html.Starter(), renderOptions},
		{notebook_server_http.Starter(), nil},

		// action managers
		{server_http_jschmhr.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		{Starter(), nil},
	}

	return starters, nil
}

func ClientComponents() ([]starter.Starter, error) {
	_, pagesConfig, restConfig, err := CompleteServerConfigs()
	if err != nil {
		return nil, err
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
		{auth_http.Starter(), common.Map{"server_config": *restConfig}}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		{files_http.Starter(), endpointsOptions},
		{records_http.Starter(), endpointsOptions},
	}

	return starters, nil
}
