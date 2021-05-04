package catalogue_www

import (
	"fmt"

	nb_www_menu2 "github.com/pavlo67/tools/app_parts/nb_www_menu"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/common/common/files/files_fs"
	"github.com/pavlo67/tools/common/actor_www"
	"github.com/pavlo67/tools/common/kv"
	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/data/entities/items/catalogue_files"
	"github.com/pavlo67/tools/components/catalogue/catalogue_server_http"
)

var _ actor_www.Operator = &catalogueActor{}

// var key = logger.GetCallInfo().PackageName

func Actor(modifyMenu thread.KVAdd, config actor_www.ConfigPages) actor_www.Operator {
	modifyMenu.Add(kv.Item{
		Key: []string{config.Prefix},
		Value: nb_www_menu2.MenuItemWWW{
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
	actorConfig actor_www.ConfigPages
	modifyMenu  thread.KVAdd
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

func (ca *catalogueActor) Config() (*actor_www.ConfigPages, error) {
	if ca == nil {
		return nil, fmt.Errorf("catalogueActor == nil")
	}

	return &ca.actorConfig, nil
}
