package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/apps/st_www/st_settings"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

const serviceName = "demo"

func main() {
	versionOnly, envPath, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, serviceName, apps.AppsSubpathDefault)
	if versionOnly {
		return
	}

	// running starters

	label := "BACKUP/SQLITE/REST BUILD"
	joinerOp, err := starter.Run(st_settings.Components(envPath, true, false), cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	st_settings.WG.Wait()
}
