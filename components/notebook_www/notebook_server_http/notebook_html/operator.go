package notebook_html

import (
	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/data_exchange/components/tags"
	"github.com/pavlo67/tools/entities/records"
)

type Operator interface {

	// complete page fragments ---------------------------------

	FragmentsView(r *records.Item, children []records.Item, message string, identity *auth.Identity) (server_http.Fragments, error)
	FragmentsEdit(r *records.Item, children []records.Item, message string, identity *auth.Identity) (server_http.Fragments, error)
	FragmentsListTagged(tag tags.Item, tagged []records.Item, identity *auth.Identity) (server_http.Fragments, error)

	// page elements  ------------------------------------------

	HTMLIndex(identity *auth.Identity) string
	HTMLFiles(recordItems []records.Item, identity *auth.Identity) string
	HTMLTags(tsm tags.StatMap, identity *auth.Identity) string
}
