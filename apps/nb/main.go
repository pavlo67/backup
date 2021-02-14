package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/apps/nb/nb_api"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

const serviceName = "demo"

func main() {
	versionOnly, _, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, serviceName, apps.AppsSubpathDefault)
	if versionOnly {
		return
	}

	starters, err := nb_api.Components(true)
	if err != nil {
		l.Fatal(err)
	}

	label := "NB/HTML/REST BUILD"
	joinerOp, err := starter.Run(starters, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	nb_api.WG.Wait()
}
