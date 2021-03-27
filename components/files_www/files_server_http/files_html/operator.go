package files_html

import (
	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"

	"github.com/pavlo67/tools/entities/files"
)

type Operator interface {

	// complete pages ------------------------------------------

	FragmentsList(filesItems []files.Item, path string, identity *auth.Identity) (server_http.Fragments, error)
	//FragmentsView(r *files.Item, path string, identity *auth.Identity) (server_http.Fragments, error)
	//FragmentsEdit(r *files.Item, path string, identity *auth.Identity) (server_http.Fragments, error)

	// page elements  ------------------------------------------

	HTMLFiles(filesItems []files.Item, identity *auth.Identity) string
}
