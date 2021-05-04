package main

import (
	"flag"
	"log"
	"sync"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/config"

	"github.com/pavlo67/tools/common/actor_www"

	"github.com/pavlo67/tools/app_parts/app_www_layout"
	"github.com/pavlo67/tools/app_parts/notebook_www"
)

var BuildDate, BuildTag, BuildCommit string
var versionOnly bool

func main() {
	log.Printf("builded: %s, tag: %s, commit: %s\n", BuildDate, BuildTag, BuildCommit)
	flag.BoolVar(&versionOnly, "v", false, "show build vars only")
	flag.Parse()

	if versionOnly {
		return
	}

	// common environment preparation --------------------------------------------------------------

	envPath, cfgService, l := apps.Prepare("_environments/")
	if cfgService == nil {
		l.Fatalf(`on apps.Run("_environments/") got nil cfgService`)
	}

	// server preparation --------------------------------------------------------------------------

	srvOp, commonChannel := app_www_layout.Run(*cfgService, l)

	// starting actors -----------------------------------------------------------------------------

	actorsWWW := []actor_www.Actor{
		{"nb", "/nb", "нотатник", notebook_www.Actor()},

		//catalogue_www2.Actor(commonChannel, actorConfigs["catalogue_home"]),
		//catalogue_www2.Actor(commonChannel, actorConfigs["catalogue_cinnamon"]),
	}

	var err error

	for _, actorWWW := range actorsWWW {
		cfgServicePath := envPath + actorWWW.Key + ".yaml"
		cfgService, err = config.Get(cfgServicePath, config.MarshalerYAML)
		if err != nil || cfgService == nil {
			l.Fatalf("on config.Get(%s, serializer.MarshalerYAML)", cfgServicePath, cfgService, err)
		}

		joinerOp, configPages, err := actorWWW.Run(*cfgService, l, actorWWW.Prefix, actor_www.Config{
			Key:      actorWWW.Key,
			Title:    actorWWW.Title,
			Callback: commonChannel,
		})
		if err != nil {
			l.Fatal(err)
		}
		defer joinerOp.CloseAll()

		if configPages != nil {
			if err := configPages.HandlePages(srvOp, l); err != nil {
				l.Fatal(err)
			}
		}
	}

	// starting http server ---------------------------------------------------------------

	var WG sync.WaitGroup
	WG.Add(1)

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	WG.Wait()
}
