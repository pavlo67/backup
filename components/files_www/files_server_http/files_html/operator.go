package files_html

import (
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/tools/entities/records"
)

type Operator interface {

	// complete pages ------------------------------------------

	CommonPage(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) (map[string]string, error)
	View(r *records.Item, children []records.Item, message string, identity *auth.Identity) (map[string]string, error)
	Edit(r *records.Item, children []records.Item, message string, identity *auth.Identity) (map[string]string, error)

	// page elements  ------------------------------------------

	HTMLFiles(recordItems []records.Item, identity *auth.Identity) string
}
