package views_html

import "github.com/pavlo67/common/common"

type Field struct {
	Key        string
	Label      string
	Type       string
	Options    common.Map
	Attributes map[string]string
}
