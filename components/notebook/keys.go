package notebook

import (
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server/server_http"
)

const InterfaceKey joiner.InterfaceKey = "notebook_www"

//const InterfaceType interfaces.Type = "notebook"

const IntefaceKeyHTMLRoot server_http.EndpointKey = "notebook_html_root"
const IntefaceKeyHTMLView server_http.EndpointKey = "notebook_html_view"
const IntefaceKeyHTMLCreate server_http.EndpointKey = "notebook_html_create"
const IntefaceKeyHTMLEdit server_http.EndpointKey = "notebook_html_edit"
const IntefaceKeyHTMLSave server_http.EndpointKey = "notebook_html_save"
const IntefaceKeyHTMLDelete server_http.EndpointKey = "notebook_html_delete"

const IntefaceKeyHTMLTags server_http.EndpointKey = "notebook_html_tags"
const IntefaceKeyHTMLTagged server_http.EndpointKey = "notebook_html_tagged"

//const IntefaceKeyRESTRead joiner.InterfaceKey = "notebook_rest_read"
//const IntefaceKeyRESTSave joiner.InterfaceKey = "notebook_rest_save"
//const IntefaceKeyRESTDele joiner.InterfaceKey = "notebook_rest_dele"
//const IntefaceKeyRESTTags joiner.InterfaceKey = "notebook_rest_tags"
//const IntefaceKeyRESTList joiner.InterfaceKey = "notebook_rest_list"
