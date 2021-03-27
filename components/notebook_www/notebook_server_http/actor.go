package notebook_server_http

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/common/actor"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/entities/files/files_fs"
	"github.com/pavlo67/tools/entities/records/records_sqlite"
)

var _ actor.OperatorWWW = &notebookActor{}

func Actor() actor.OperatorWWW {
	return &notebookActor{}
}

type notebookActor struct {
}

func (*notebookActor) Name() string {
	return ""
}

var filesOptions = common.Map{
	"base_path": "../_files_fs_test",
}

func (*notebookActor) Starters(options common.Map) ([]starter.Starter, error) {
	starters := []starter.Starter{
		// general purposes components
		{db_sqlite.Starter(), nil},

		//// auth/persons components
		//{persons_fs.Starter(), nil},
		//{auth_persons.Starter(), nil},
		//{auth_jwt.Starter(), nil},
		//{auth_server_http.Starter(), nil}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		{files_fs.Starter(), filesOptions},
		{records_sqlite.Starter(), nil},
		{Starter(), nil},

		// action managers

		// actions starter (connecting specific actions to the corresponding action managers)
		// {nb_www_settings.Starter(), nil},
	}

	return starters, nil
}

func (*notebookActor) Config() (server_http.ConfigPages, error) {
	return PagesConfig, nil
}

//func ClientComponents() ([]starter.Starter, error) {
//
//	if err := EndpointsPageSettled.CompleteDirectly(notebook_server_http.EndpointsSettled, "", 0, pagesPrefix); err != nil {
//		return nil, err
//	}
//
//	//if err := RestConfig.CompleteDirectly(auth_server_http.Config, "", 0, restPrefix); err != nil {
//	//	return nil, err
//	//}
//
//	endpointsOptions := common.Map{
//		"pages_config": &EndpointsPageSettled,
//		// "rest_config":  restConfig,
//	}
//
//	starters := []starter.Starter{
//		// general purposes components
//		{control.Starter(), nil},
//
//		// auth/persons components
//		{auth_jwt.Starter(), nil},
//		{auth_http.Starter(), common.Map{"server_config": RestConfig}}, // common.Map{"auth_jwt_key": ""}
//
//		// notebook components
//		{files_http.Starter(), endpointsOptions},
//		{records_http.Starter(), endpointsOptions},
//	}
//
//	return starters, nil
//}
