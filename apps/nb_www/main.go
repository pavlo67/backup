package main

import (
	"io/ioutil"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/tools/common/actor"

	"github.com/pavlo67/tools/components/notebook_www/notebook_server_http"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

func main() {
	versionOnly, _, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, "_environments/")
	if versionOnly {
		return
	}

	templatePath := filelib.CurrentPath() + "templates/local.html"
	htmlTemplateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		l.Fatalf("on ioutil.ReadFile(%s): %s", templatePath, err)
	}

	actorsWWW := []actor.OperatorWWW{
		notebook_server_http.Actor(),
	}

	joinerOps, err := actor.RunWWW(cfgService, string(htmlTemplateBytes), "NB/HTML/REST BUILD", actorsWWW, l)
	for _, joinerOp := range joinerOps {
		defer joinerOp.CloseAll()
	}

	l.Fatal(err)

}
