package tags_html

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/tools/components/formatter"
	"github.com/pavlo67/tools/components/tags"
)

const List formatter.Key = "tags_list"

// should be thread-safe
type Operator interface {
	Prepare(key formatter.Key, template string, params common.Map) error

	HTMLTags(rs tags.Stats) (string, error)
}
