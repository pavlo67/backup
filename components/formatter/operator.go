package formatter

import "github.com/pavlo67/common/common"

type Key string

// should be thread-safe
type Operator interface {
	Prepare(key Key, template string, params common.Map) error
	Format(value interface{}, key Key) (string, error)
}
