package formatter_records_html

import (
	"encoding/json"

	"github.com/pavlo67/tools/components/records"
)

func Edit(item *records.Item) string {
	bytes, _ := json.Marshal(item)

	return string(bytes)
}
