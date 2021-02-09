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
)

func Components(startREST bool) []starter.Starter {

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},

		// auth/persons components
		{persons_fs.Starter(), nil},
		{auth_persons.Starter(), nil},
		{auth_jwt.Starter(), common.Map{"interface_key": auth_jwt.InterfaceKey}},
		{auth_server_http.Starter(), nil},
	}

	if !startREST {
		return starters
	}

	starters = append(
		starters,

		// action managers
		starter.Starter{server_http_jschmhr.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		starter.Starter{Starter(), nil},
	)

	return starters
}
