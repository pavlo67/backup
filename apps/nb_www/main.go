package main

import (
	"io/ioutil"
	"os"

	"github.com/pavlo67/tools/components/catalogue/catalogue_www"

	"github.com/pavlo67/tools/components/notebook/notebook_www"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/tools/common/actor"
	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/tools/apps/nb_www/nb_www_menu"
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

	processMenu, err := thread.NewFIFOKVItems(&nb_www_menu.MenuWWW{})
	if err != nil {
		l.Fatalf("on thread.NewFIFOKVItems(): %s", err)
	}

	var actorConfigs map[string]actor.Config
	if err = cfgService.Value("actors", &actorConfigs); err != nil {
		l.Fatalf(`on cfgService.Value("actors", &actorConfigs): %s`, err)
	}

	//l.Infof("%#v", actorConfigs["catalogue_home"])
	//l.Fatalf("%#v", actorConfigs["catalogue_cinnamon"])

	actorsWWW := []actor.OperatorWWW{
		notebook_www.Actor(processMenu, actorConfigs["notebook"]),
		catalogue_www.Actor(processMenu, actorConfigs["catalogue_home"]),
		catalogue_www.Actor(processMenu, actorConfigs["catalogue_cinnamon"]),
	}

	joinerOps, err := actor.RunWWW(
		cfgService, "NB/HTML/REST BUILD",
		string(htmlTemplateBytes), staticPath, processMenu,
		actorsWWW,
		l,
	)
	for _, joinerOp := range joinerOps {
		if joinerOp != nil {
			defer joinerOp.CloseAll()
		}
	}

	l.Fatal(err)

}
