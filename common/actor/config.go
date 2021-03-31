package actor

import "github.com/pavlo67/common/common"

type Config struct {
	Type    string                `yaml:"type"`
	Title   string                `yaml:"title"`
	Prefix  string                `yaml:"prefix"`
	Order   int                   `yaml:"order"`
	Options map[string]common.Map `yaml:"options"`
}
