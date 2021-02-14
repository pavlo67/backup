package records_html

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/tools/components/tags"

	"github.com/pavlo67/tools/components/formatter"
	"github.com/pavlo67/tools/components/records"
)

const Full formatter.Key = "full"
const Brief formatter.Key = "brief"
const Edit formatter.Key = "edit"
const Tag formatter.Key = "tag"

// should be thread-safe
type Operator interface {
	Prepare(key formatter.Key, template string, params common.Map) error

	HTMLView(r *records.Item, children []records.Item) (string, error)
	HTMLTagged(tag tags.Item, tagged []records.Item) (string, error)
}
