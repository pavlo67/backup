package server_http

import (
	"strconv"

	"github.com/pavlo67/common/common/errors"
)

type ConfigCommon struct {
	Title   string
	Version string
	Host    string
	Port    string
	Prefix  string
}

func (c *ConfigCommon) Complete(host string, port int, prefix string) error {
	if c == nil {
		return errors.New("no server_http.Config to complete")
	}

	var portStr string
	if port > 0 {
		portStr = ":" + strconv.Itoa(port)
	}

	c.Host, c.Port, c.Prefix = host, portStr, prefix

	return nil
	//for key, ep := range c.Config {
	//	if endpoint, ok := joinerOp.Interface(key).(EndpointREST); ok {
	//		ep.EndpointREST = endpoint
	//		c.Config[key] = ep
	//	} else if endpointPtr, _ := joinerOp.Interface(key).(*EndpointREST); endpointPtr != nil {
	//		ep.EndpointREST = *endpointPtr
	//		c.Config[key] = ep
	//	} else {
	//		return fmt.Errorf("no server_http.EndpointREST joined with key %s", key)
	//	}
	//}

}
