package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/apps/nb_www/nb_www_settings"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

func main() {
	versionOnly, _, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, "../_environments/")
	if versionOnly {
		return
	}

	starters, err := nb_www_settings.ServerComponents()
	if err != nil {
		l.Fatal(err)
	}

	label := "NB/HTML/REST BUILD"
	joinerOp, err := starter.Run(starters, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	nb_www_settings.WG.Wait()
}
