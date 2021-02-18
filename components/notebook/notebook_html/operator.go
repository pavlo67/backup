package notebook_html

import (
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

type Key string

const Full Key = "full"
const Brief Key = "brief"
const Edit Key = "edit"
const Tag Key = "tag"

// should be thread-safe
type Operator interface {
	// DEPRECATED
	HTMLPage(title, htmlHeader, htmlIndex, htmlContent, htmlMessage string) server.Response

	HTMLError(httpStatus int, htmlError string) (server.Response, error)

	HTMLRoot(htmlHello, htmlTags string) (server.Response, error)

	HTMLView(r *records.Item, children []records.Item, message string) (server.Response, error)
	HTMLEdit(r *records.Item, children []records.Item, message string) (server.Response, error)
	HTMLList(tag tags.Item, tagged []records.Item) (string, error)
	HTMLTags(rs tags.Stats) (string, error)

	HTMLMessage(errs errors.Error) string
}

// Prepare(key Key, template string, params common.Map) error
