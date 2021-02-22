package notebook_html

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"

	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

type Operator interface {

	// complete pages ------------------------------------------

	CommonPage(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) (string, error)

	View(r *records.Item, children []records.Item, message string, options *crud.Options) (string, error)
	Edit(r *records.Item, children []records.Item, message string, options *crud.Options) (string, error)
	ListTagged(tag tags.Item, tagged []records.Item, options *crud.Options) (string, error)

	// page elements  ------------------------------------------

	HTMLIndex(options *crud.Options) string
	HTMLRecords(recordItems []records.Item, identity *auth.Identity) string
	HTMLTags(tsm tags.StatMap, options *crud.Options) string
}
