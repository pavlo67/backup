package notebook_www

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/common/actor"
	"github.com/pavlo67/tools/common/kv"
	"github.com/pavlo67/tools/common/thread"
	"github.com/pavlo67/tools/components/notebook/notebook_server_http"

	"github.com/pavlo67/tools/entities/records/records_sqlite"

	"github.com/pavlo67/tools/apps/nb_www/nb_www_menu"
)

var _ actor.OperatorWWW = &notebookActor{}

// var key = logger.GetCallInfo().PackageName

func Actor(modifyMenu thread.FIFOKVItemsAdd, config actor.Config) actor.OperatorWWW {
	modifyMenu.Add(kv.Item{
		Key: []string{config.Prefix},
		Value: nb_www_menu.MenuItemWWW{
			HRef:  "/" + config.Prefix,
			Title: config.Title,
		},
	})

	return &notebookActor{
		actorConfig: config,
		modifyMenu:  modifyMenu,
	}
}

type notebookActor struct {
	actorConfig actor.Config
	modifyMenu  thread.FIFOKVItemsAdd
}

func (na *notebookActor) Name() string {
	if na == nil {
		return ""
	}

	return na.actorConfig.Title
}

func (na *notebookActor) Starters() ([]starter.Starter, error) {
	if na == nil {
		return nil, fmt.Errorf("notebookActor == nil")
	}

	starters := []starter.Starter{
		// general purposes components
		{db_sqlite.Starter(), nil},

		//// auth/persons components
		//{persons_fs.Starter(), nil},
		//{auth_persons.Starter(), nil},
		//{auth_jwt.Starter(), nil},
		//{auth_server_http.Starter(), nil}, // common.Map{"auth_jwt_key": ""}

		// notebook components
		// {files_fs.Starter(), filesOptions},
		{records_sqlite.Starter(), nil},
		{notebook_server_http.Starter(), common.Map{"prefix": na.actorConfig.Prefix}},

		// action managers

		// actions starter (connecting specific actions to the corresponding action managers)
		// {nb_www_settings.Starter(), nil},
	}

	return starters, nil
}

func (na *notebookActor) Config() (*actor.Config, error) {
	if na == nil {
		return nil, fmt.Errorf("notebookActor == nil")
	}

	return &na.actorConfig, nil
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
