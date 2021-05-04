package notebook_www

import (
	"fmt"

	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/tools/app_parts/app_www_layout"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records/records_sqlite"

	"github.com/pavlo67/tools/common/actor_www"
	"github.com/pavlo67/tools/common/kv"

	"github.com/pavlo67/tools/components/notebook/notebook_server_http"
)

var _ actor_www.Operator = &notebookActor{}

// var key = logger.GetCallInfo().PackageName

func Actor() actor_www.Operator {
	return &notebookActor{}
}

type notebookActor struct {
	//prefix string
	//config actor_www.Config
}

const onRun = "on notebookActor.Run()"

func (na *notebookActor) Run(cfgService config.Config, l logger.Operator, prefix string, config actor_www.Config) (joiner.Operator, *actor_www.ConfigPages, error) {
	if na == nil {
		return nil, nil, fmt.Errorf("notebookActor == nil")
	}

	if config.Callback == nil {
		return nil, nil, errors.New(onRun + ": no config.Callback is defined")
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
		{notebook_server_http.Starter(), common.Map{"prefix": prefix}},
	}

	joinerOp, err := starter.Run(starters, &cfgService, "NOTEBOOK_WWW_ACTOR BUILD", l)
	if err != nil {
		return nil, nil, fmt.Errorf(onRun+": got %s", err)
	}

	var configPages *actor_www.ConfigPages

	switch v := joinerOp.Interface(actor_www.ConfigPagesInterfaceKey).(type) {
	case actor_www.ConfigPages:
		configPages = &v
	case *actor_www.ConfigPages:
		configPages = v
	}

	if configPages == nil {
		return nil, nil, errors.New(onRun + ": no ConfigPages is got")

	}

	config.Callback.Add(kv.Item{
		Key: []string{prefix},
		Value: app_www_layout.MenuItemWWW{
			HRef:  "/" + prefix,
			Title: config.Title,
		},
	})

	return joinerOp, configPages, nil
}
