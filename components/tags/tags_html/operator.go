package tags_html

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/tools/components/tags"
)

type Key string

const List Key = "tags_list"

// should be thread-safe
type Operator interface {
	Prepare(key Key, template string, params common.Map) error

	HTMLTags(rs tags.Stats) (string, error)
}
