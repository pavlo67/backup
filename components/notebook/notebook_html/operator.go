package notebook_html

import (
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data_exchange/components/tags"
	"github.com/pavlo67/tools/components/records"
)

type Operator interface {

	// complete pages ------------------------------------------

	CommonPage(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) (string, error)

	View(r *records.Item, children []records.Item, message string, identity *auth.Identity) (string, error)
	Edit(r *records.Item, children []records.Item, message string, identity *auth.Identity) (string, error)
	ListTagged(tag tags.Item, tagged []records.Item, identity *auth.Identity) (string, error)

	// page elements  ------------------------------------------

	HTMLIndex(identity *auth.Identity) string
	HTMLRecords(recordItems []records.Item, identity *auth.Identity) string
	HTMLTags(tsm tags.StatMap, identity *auth.Identity) string
}
