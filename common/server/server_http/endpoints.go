package server_http

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server/server_http"
)

// this trick allows to prevent run-time errors with wrong endpoint parameters number
// using CheckGet...() functions we move parameter number checks to initiation stage

type Get1 func(string) (string, error)
type Get2 func(string, string) (string, error)
type Get3 func(string, string, string) (string, error)
type Get4 func(string, string, string, string) (string, error)

func CheckGet0(c server_http.Config, endpointKey joiner.InterfaceKey, createFullURL bool) (string, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if strings.ToUpper(ep.Method) != "GET" {
		return "", fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	return urlStr, nil
}

func CheckGet1(c server_http.Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get1, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1 string) (string, error) {
		p1 = strings.TrimSpace(p1)
		if p1 == "" {
			return "", fmt.Errorf("empty param %s for endpoint (%s / %#v)", p1, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(p1)
		return urlStr, nil
	}, nil
}

func CheckGet2(c server_http.Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get2, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2 string) (string, error) {
		params := [2]string{p1, p2}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet3(c server_http.Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get3, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3 string) (string, error) {
		params := [3]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet4(c server_http.Config, endpointKey joiner.InterfaceKey, createFullURL bool) (Get4, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3, p4 string) (string, error) {
		params := [4]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}
