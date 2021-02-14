package records_html

import (
	"encoding/json"

	"github.com/pavlo67/tools/components/records"
)

func HTMLEdit(item *records.Item) string {
	bytes, _ := json.Marshal(item)

	return string(bytes)
}
