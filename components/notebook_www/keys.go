package notebook_www

import (
	"github.com/pavlo67/common/common/joiner"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

const InterfaceKey joiner.InterfaceKey = "notebook"

const IntefaceKeyHTMLRoot server_http.EndpointKey = "notebook_html_root"

const IntefaceKeyHTMLView server_http.EndpointKey = "notebook_html_view"
const IntefaceKeyHTMLCreate server_http.EndpointKey = "notebook_html_create"
const IntefaceKeyHTMLEdit server_http.EndpointKey = "notebook_html_edit"
const IntefaceKeyHTMLSave server_http.EndpointKey = "notebook_html_save"
const IntefaceKeyHTMLDelete server_http.EndpointKey = "notebook_html_delete"

const IntefaceKeyHTMLTags server_http.EndpointKey = "notebook_html_tags"
const IntefaceKeyHTMLTagged server_http.EndpointKey = "notebook_html_tagged"

//
//const IntefaceKeyRESTRead joiner.InterfaceKey = "notebook_rest_read"
//const IntefaceKeyRESTSave joiner.InterfaceKey = "notebook_rest_save"
//const IntefaceKeyRESTDele joiner.InterfaceKey = "notebook_rest_dele"
//const IntefaceKeyRESTTags joiner.InterfaceKey = "notebook_rest_tags"
//const IntefaceKeyRESTList joiner.InterfaceKey = "notebook_rest_list"
