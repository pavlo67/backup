package server_http

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server/server_http"
)

func EP(serverConfig server_http.Config, endpointKey joiner.InterfaceKey, params []string, createFullURL bool) (string, string, error) {
	ep, ok := serverConfig.EndpointsSettled[endpointKey]
	if !ok {
		return "", "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if len(ep.PathParams) != len(params) {
		return "", "", fmt.Errorf("wrong params list (%#v) for endpoint (%s / %#v)", params, endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = serverConfig.Host
		if serverConfig.Port = strings.TrimSpace(serverConfig.Port); serverConfig.Port != "" {
			urlStr += ":" + serverConfig.Port
		}
	}
	urlStr += serverConfig.Prefix

	for i, param := range params {
		if param == "" {
			return "", "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(param)
	}

	return ep.Method, urlStr, nil
}
