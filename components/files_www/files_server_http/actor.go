package files_server_http

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/apps/nb_www/nb_www_menu"
	"github.com/pavlo67/tools/common/kv"
	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/tools/common/actor"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/tools/entities/files/files_fs"
)

var _ actor.OperatorWWW = &filesActor{}

var key = logger.GetCallInfo().PackageName

func Actor(modifyMenu thread.FIFOKVItemsAdd, prefix string) actor.OperatorWWW {
	modifyMenu.Add(kv.Item{
		Key: []string{key},
		Value: nb_www_menu.MenuItemWWW{
			HRef:  "/" + prefix + "/list",
			Title: key,
		},
	})

	return &filesActor{
		modifyMenu: modifyMenu,
	}
}

type filesActor struct {
	modifyMenu thread.FIFOKVItemsAdd
}

func (*filesActor) Name() string {
	return key
}

var filesOptions = common.Map{
	"base_path": "../_files_fs_test",
}

func (*filesActor) Starters(options common.Map) ([]starter.Starter, error) {
	starters := []starter.Starter{
		{files_fs.Starter(), filesOptions},
		{Starter(), nil},
	}

	return starters, nil
}

func (*filesActor) Config() (server_http.ConfigPages, error) {
	return PagesConfig, nil
}
