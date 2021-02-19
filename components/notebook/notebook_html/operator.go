package notebook_html

import (
	"net/http"

	"github.com/pavlo67/common/common/server"

	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
)

//type Key string
//
//const Full Key = "full"
//const Brief Key = "brief"
//const Edit Key = "edit"
//const Tag Key = "tag"

// should be thread-safe
type Operator interface {
	HTMLRoot(htmlHello string, tagsStatMap tags.StatMap) (server.Response, error)
	HTMLError(httpStatus int, err error, publicDetails string, req *http.Request) (server.Response, error)

	HTMLView(r *records.Item, children []records.Item, message string) (server.Response, error)
	HTMLEdit(r *records.Item, children []records.Item, message string) (server.Response, error)

	HTMLTags(tsm tags.StatMap) (server.Response, error)
	HTMLTagged(tag tags.Item, tagged []records.Item) (server.Response, error)
}

// Prepare(key Key, template string, params common.Map) error
