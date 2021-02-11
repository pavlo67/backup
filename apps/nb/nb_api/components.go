package nb_api

import (
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

	"github.com/pavlo67/tools/components/records/formatter_records_html"
	"github.com/pavlo67/tools/components/records/records_sqlite"
)

func Components(startServer bool) []starter.Starter {

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
		{files_fs.Starter(), common.Map{"buckets": files.Buckets{files.BucketID("1"): "1"}}},
		{records_sqlite.Starter(), nil},
		{formatter_records_html.Starter(), nil},
		{notebook_server_http.Starter(), nil},
	}

	if !startServer {
		return starters
	}

	starters = append(
		starters,

		// action managers
		starter.Starter{server_http_jschmhr.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		starter.Starter{Starter(), common.Map{"prefix_rest": "/rest", "prefix_pages": ""}},
	)

	return starters
}
