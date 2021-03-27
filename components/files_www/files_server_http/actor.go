package files_server_http

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/common/actor"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/tools/components/files_www/files_server_http/files_html"
	"github.com/pavlo67/tools/entities/files/files_fs"
)

var _ actor.OperatorWWW = &filesActor{}

func Actor() actor.OperatorWWW {
	return &filesActor{}
}

type filesActor struct {
}

func (*filesActor) Name() string {
	return ""
}

var filesOptions = common.Map{
	"base_path": "../_files_fs_test",
}

func (*filesActor) Starters(options common.Map) ([]starter.Starter, error) {
	renderOptions := common.Map{
		"pages_config": &PagesConfig,
	}

	starters := []starter.Starter{
		{files_fs.Starter(), filesOptions},
		{files_html.Starter(), renderOptions},
		{Starter(), nil},
	}

	return starters, nil
}

func (*filesActor) Config() (server_http.ConfigPages, error) {
	return PagesConfig, nil
}
