package server_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/server"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/logger"
)

type WorkerHTTPREST func(OperatorV2, *http.Request, PathParams, *auth.Identity) (server.Response, error)

type EndpointREST struct {
	EndpointDescription
	WorkerHTTPREST
}

type EndpointRESTSettled struct {
	Path     string
	Tags     []string `json:",omitempty"`
	Produces []string `json:",omitempty"`
	EndpointREST
}

type EndpointsREST map[EndpointKey]EndpointRESTSettled

type ConfigREST struct {
	ConfigCommon
	Endpoints EndpointsREST
}

func (c ConfigREST) EP(endpointKey EndpointKey, params []string, createFullURL bool) (string, string, error) {
	ep, ok := c.Endpoints[endpointKey]
	if !ok {
		return "", "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if len(ep.PathParams) != len(params) {
		return "", "", fmt.Errorf("wrong params list (%#v) for endpoint (%s / %#v)", params, endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	for i, param := range params {
		if param == "" {
			return "", "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(param)
	}

	return ep.Method, urlStr, nil
}

type Swagger map[string]interface{}

func (c ConfigREST) SwaggerV2(isHTTPS bool) ([]byte, error) {
	paths := map[string]common.Map{} // map[string]map[string]map[string]interface{}{}

	for key, ep := range c.Endpoints {

		path := c.Prefix + ep.PathTemplateBraced(ep.Path)
		method := strings.ToLower(ep.Method)

		epDescr := common.Map{
			"operationId": key,
			"tags":        ep.Tags,
		}

		if len(ep.Produces) >= 1 {
			epDescr["produces"] = ep.Produces
		} else {
			epDescr["produces"] = []string{"application/json"}
		}

		var parameters []interface{} // []map[string]interface{}

		for _, pp := range ep.PathParams {
			if len(pp) > 0 && pp[0] == '*' {
				pp = pp[1:]
			}

			parameters = append(
				parameters,
				common.Map{
					"in":          "path",
					"required":    true,
					"name":        pp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}
		for _, qp := range ep.QueryParams {
			parameters = append(
				parameters,
				common.Map{
					"in":          "query",
					"required":    false, // TODO!!!
					"name":        qp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}

		if method == "post" {
			if len(ep.BodyParams) > 0 {
				parameters = append(parameters, ep.BodyParams)
			} else {
				parameters = append(parameters, common.Map{
					"in":       "body",
					"required": true,
					"name":     "body_item",
					"type":     "string",
				})
			}
		}

		if len(parameters) > 0 {
			epDescr["parameters"] = parameters
		}

		if epDescrPrev, ok := paths[path][method]; ok {
			return nil, fmt.Errorf("duplicate endpoint description (%s %s): \n%#v\nvs.\n%#v", method, path, epDescrPrev, epDescr)
		}
		if _, ok := paths[path]; ok { // pathPrev
			paths[path][method] = epDescr
		} else {
			paths[path] = common.Map{method: epDescr} // map[string]map[string]interface{}
		}
	}

	var schemes []string
	if isHTTPS {
		schemes = []string{"https", "http"}
	} else {
		schemes = []string{"http"}
	}

	swagger := Swagger{
		"swagger": "2.0",
		"info": map[string]string{
			"title":   c.Title,
			"version": c.Version,
		},
		// "basePath": c.Prefix,
		"schemes": schemes,
		"port":    c.Port,
		"paths":   paths,
	}

	return json.MarshalIndent(swagger, "", " ")
}

func (c ConfigREST) InitSwagger(isHTTPS bool, swaggerStaticFilePath string, l logger.Operator) error {
	swaggerJSON, err := c.SwaggerV2(isHTTPS)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(swaggerStaticFilePath, swaggerJSON, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerStaticFilePath, swaggerJSON, err)
	}
	l.Infof("%d bytes are written into %s", len(swaggerJSON), swaggerStaticFilePath)

	return nil
}

const onHandleEndpoints = "on server_http.HandleEndpoints()"

func (c ConfigREST) HandleEndpoints(srvOp OperatorV2, l logger.Operator) error {
	if srvOp == nil {
		return errors.New(onHandleEndpoints + ": srvOp == nil")
	}

	for key, ep := range c.Endpoints {
		if err := srvOp.Handle(key, c.Prefix+ep.Path, WrapperHTTPREST, ep.EndpointREST); err != nil {
			return fmt.Errorf(onHandleEndpoints+": handling %s, %s, %#v got %s", key, ep.Path, ep, err)
		}
	}

	return nil
}
