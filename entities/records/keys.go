package records

import (
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "records"
const InterfaceCleanerKey joiner.InterfaceKey = "records_cleaner"

const CollectionDefault = "records"

const HasTag selectors.Key = "has_tag"
const HasNoTag selectors.Key = "has_no_tag"
const HasParent selectors.Key = "has_parent"
