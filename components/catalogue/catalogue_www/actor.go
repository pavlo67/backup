package catalogue_www

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/common/actor"
	"github.com/pavlo67/tools/common/files/files_fs"
	"github.com/pavlo67/tools/common/kv"
	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/tools/components/catalogue/catalogue_server_http"
	"github.com/pavlo67/tools/entities/items/catalogue_files"

	"github.com/pavlo67/tools/apps/nb_www/nb_www_menu"
)

var _ actor.OperatorWWW = &catalogueActor{}

// var key = logger.GetCallInfo().PackageName

func Actor(modifyMenu thread.FIFOKVItemsAdd, config actor.Config) actor.OperatorWWW {
	modifyMenu.Add(kv.Item{
		Key: []string{config.Prefix},
		Value: nb_www_menu.MenuItemWWW{
			HRef:  "/" + config.Prefix + "/list",
			Title: config.Title,
		},
	})

	return &catalogueActor{
		actorConfig: config,
		modifyMenu:  modifyMenu,
	}
}

type catalogueActor struct {
	actorConfig actor.Config
	modifyMenu  thread.FIFOKVItemsAdd
}

func (ca *catalogueActor) Name() string {
	if ca == nil {
		return ""
	}
	return ca.actorConfig.Title
}

func (ca *catalogueActor) Starters() ([]starter.Starter, error) {
	if ca == nil {
		return nil, fmt.Errorf("catalogueActor == nil")
	}

	filesFSConfigKey := ca.actorConfig.Options["files_fs"].StringDefault("config_key", "")

	//log.Printf("%s --> %#v --> %s", ca.actorConfig.Prefix, ca.actorConfig.Options, filesFSConfigKey)

	starters := []starter.Starter{
		{files_fs.Starter(), common.Map{"config_key": filesFSConfigKey}},
		{catalogue_files.Starter(), nil},
		{catalogue_server_http.Starter(), common.Map{"prefix": ca.actorConfig.Prefix}}}

	return starters, nil
}

func (ca *catalogueActor) Config() (*actor.Config, error) {
	if ca == nil {
		return nil, fmt.Errorf("catalogueActor == nil")
	}

	return &ca.actorConfig, nil
}
