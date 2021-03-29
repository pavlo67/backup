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

const title = "файли"

func Actor(modifyMenu thread.FIFOKVItemsAdd, options common.Map) actor.OperatorWWW {
	prefix := options.StringDefault("prefix", "")
	modifyMenu.Add(kv.Item{
		Key: []string{key},
		Value: nb_www_menu.MenuItemWWW{
			HRef:  "/" + prefix + "/list",
			Title: title,
		},
	})

	return &filesActor{
		options:    options,
		modifyMenu: modifyMenu,
	}
}

type filesActor struct {
	options    common.Map
	modifyMenu thread.FIFOKVItemsAdd
}

func (*filesActor) Name() string {
	return title
}

func (fa *filesActor) Options() common.Map {
	if fa == nil {
		return nil
	}
	return fa.options
}

var filesOptions = common.Map{
	"base_path": "../_files_fs_test",
}

func (fa *filesActor) Starters() ([]starter.Starter, error) {
	prefix := fa.Options().StringDefault("prefix", "")

	starters := []starter.Starter{
		{files_fs.Starter(), filesOptions},
		{Starter(), common.Map{"prefix": prefix}}}

	return starters, nil
}

func (*filesActor) Config() (*server_http.ConfigPages, error) {
	return &PagesConfig, nil
}
