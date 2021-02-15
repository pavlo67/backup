package records_html

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/tools/components/tags"

	"github.com/pavlo67/tools/components/records"
)

type Key string

const Full Key = "full"
const Brief Key = "brief"
const Edit Key = "edit"
const Tag Key = "tag"

// should be thread-safe
type Operator interface {
	Prepare(key Key, template string, params common.Map) error

	HTMLView(r *records.Item, children []records.Item) (string, error)
	HTMLTagged(tag tags.Item, tagged []records.Item) (string, error)
}
