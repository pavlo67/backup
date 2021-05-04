package actor_www

import (
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"

	"github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/common/thread"
)

type ConfigPages = server_http_v2.ConfigPages

type Config struct {
	//Name string
	Key      string
	Title    string
	Callback thread.KVAdd
}

type Actor struct {
	Key    string
	Prefix string
	Title  string
	Operator
}

type Operator interface {
	// Name() string
	Run(cfgService config.Config, l logger.Operator, Prefix string, config Config) (joiner.Operator, *ConfigPages, error)
}
