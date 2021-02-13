package formatter

import "github.com/pavlo67/common/common"

type Key string

const Full Key = "full"
const Brief Key = "brief"
const Edit Key = "edit"
const Tag Key = "tag"

type Operator interface {
	// should be thread-safe
	Prepare(key Key, template string, params common.Map) error

	// should be thread-safe
	Format(value interface{}, key Key) (string, error)
}
