package formatter_records_html

import (
	"encoding/json"

	"github.com/pavlo67/tools/components/records"
)

func Full(item *records.Item) string {
	bytes, _ := json.Marshal(item)

	return string(bytes)
}

func Brief(item *records.Item) string {
	bytes, _ := json.Marshal(item)

	return string(bytes)
}
