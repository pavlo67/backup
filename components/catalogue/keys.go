package catalogue

import (
	"github.com/pavlo67/common/common/joiner"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

const InterfaceKey joiner.InterfaceKey = "catalogue_www"

//const InterfaceType interfaces.Type = "catalogue"

const IntefaceKeyHTMLView server_http.EndpointKey = "catalogue_html_view"
const IntefaceKeyHTMLList server_http.EndpointKey = "catalogue_html_list"
const IntefaceKeyHTMLCreate server_http.EndpointKey = "catalogue_html_create"
const IntefaceKeyHTMLEdit server_http.EndpointKey = "catalogue_html_edit"
const IntefaceKeyHTMLSave server_http.EndpointKey = "catalogue_html_save"
const IntefaceKeyHTMLDelete server_http.EndpointKey = "catalogue_html_delete"
