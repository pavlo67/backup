package server_http

import (
	"encoding/json"
	"regexp"
	"strings"
)

type EndpointKey string

type EndpointDescription struct {
	Method      string          `json:",omitempty"`
	PathParams  []string        `json:",omitempty"`
	QueryParams []string        `json:",omitempty"`
	BodyParams  json.RawMessage `json:",omitempty"`
}

var rePathParam = regexp.MustCompile(":[^/]+")

func (ed EndpointDescription) PathTemplate(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ed.PathParams) < 1 {
		return serverPath
	}

	var pathParams []string
	for _, pp := range ed.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp)
		} else {
			pathParams = append(pathParams, ":"+pp)
		}

	}

	return serverPath + "/" + strings.Join(pathParams, "/")
}

func (ed EndpointDescription) PathTemplateBraced(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ed.PathParams) < 1 {
		return serverPath
	}

	var pathParams []string
	for _, pp := range ed.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp[1:])
		} else {
			pathParams = append(pathParams, pp)
		}

	}

	return serverPath + "/{" + strings.Join(ed.PathParams, "}/{") + "}"
}

//func (ep Endpoint) PathWithParams(params ...string) string {
//	matches := rePathParam.FindAllStringSubmatchIndex(ep.Path, -1)
//
//	numMatches := len(matches)
//	if len(params) < numMatches {
//		numMatches = len(params)
//	}
//
//	path := ep.Path
//	for nm := numMatches - 1; nm >= 0; nm-- {
//		path = path[:matches[nm][0]] + url.PathEscape(strings.ReplaceTags(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
//	}
//
//	return path
//}
