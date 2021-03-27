package main

import (
	"io/ioutil"
	"os"

	"github.com/pavlo67/tools/components/files_www/files_server_http"

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

	// static files & templates preparation --------------------------------------------------------

	templatePath := filelib.CurrentPath() + "templates/local.html"
	htmlTemplateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		l.Fatalf("reading template (%s): %s", templatePath, err)
	}

	staticPath := filelib.CurrentPath() + "static/"
	fileInfo, err := os.Stat(staticPath)
	if err != nil {
		l.Fatalf("checking static files path (%s): %s", staticPath, err)
	}
	if !fileInfo.IsDir() {
		l.Fatalf("static files path (%s) is not a directory", staticPath)
	}

	// actors start --------------------------------------------------------------------------------

	actorsWWW := []actor.OperatorWWW{
		notebook_server_http.Actor(),
		files_server_http.Actor(),
	}

	joinerOps, err := actor.RunWWW(cfgService, string(htmlTemplateBytes), staticPath, "NB/HTML/REST BUILD", actorsWWW, l)
	for _, joinerOp := range joinerOps {
		defer joinerOp.CloseAll()
	}

	l.Fatal(err)

}
