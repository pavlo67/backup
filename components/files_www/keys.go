package files_www

import (
	"github.com/pavlo67/common/common/joiner"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

const InterfaceKey joiner.InterfaceKey = "files_www"

const IntefaceKeyHTMLRoot server_http.EndpointKey = "files_html_root"
const IntefaceKeyHTMLView server_http.EndpointKey = "files_html_view"
const IntefaceKeyHTMLCreate server_http.EndpointKey = "files_html_create"
const IntefaceKeyHTMLEdit server_http.EndpointKey = "files_html_edit"
const IntefaceKeyHTMLSave server_http.EndpointKey = "files_html_save"
const IntefaceKeyHTMLDelete server_http.EndpointKey = "files_html_delete"
