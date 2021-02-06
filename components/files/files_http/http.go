package files_http

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/files"
)

var _ files.Operator = &filesHTTP{}

type filesHTTP struct {
	host      string
	endpoints map[string]server_http.EndpointConfig

	logfile string
}

const onNew = "on filesHTTP.New()"

func New(access config.Access, prefix string, endpoints server_http.Endpoints, mockHandlers bool, logfile string) (files.Operator, error) {
	filesOp := filesHTTP{
		host:      access.Host,
		endpoints: map[string]server_http.EndpointConfig{},
		logfile:   logfile,
	}

	if access.Port > 0 {
		filesOp.host += ":" + strconv.Itoa(access.Port)
	}
	filesOp.host += prefix

	var ok bool

	for _, epKey := range []string{
		files.EPRead, files.EPSave, files.EPRemove, files.EPList, files.EPStat,
	} {
		filesOp.endpoints[epKey], ok = endpoints[epKey]
		if !ok {
			return nil, errors.Errorf(onNew+": no '%s' endpoint", epKey)
		} else if filesOp.endpoints[epKey].Handler == nil {
			if mockHandlers {
				ep := filesOp.endpoints[epKey]
				ep.Handler = &server_http.Endpoint{}
				switch epKey {
				case files.EPSave:
					ep.Handler.Method = "POST"
				case files.EPRemove:
					ep.Handler.Method = "DELETE"
				default:
					ep.Handler.Method = "GET"
				}
				filesOp.endpoints[epKey] = ep
			} else {
				return nil, errors.Errorf(onNew+": no '%s' endpoint.Handler", epKey)
			}
		}

	}

	return &filesOp, nil
}

func (filesOp *filesHTTP) Save(bucketID files.BucketID, path, newFilePattern string, data []byte, options *crud.Options) (string, error) {
	if path = strings.TrimSpace(path); path == "" {
		path = "."
	}
	if newFilePattern = strings.TrimSpace(newFilePattern); newFilePattern == "" {
		newFilePattern = "."
	}

	ep := filesOp.endpoints[files.EPSave]
	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + newFilePattern

	var correctedPath string
	return correctedPath, server_http.Request(serverURL, ep, nil, &correctedPath, options.GetIdentity(), filesOp.logfile)
}

func (filesOp *filesHTTP) Read(bucketID files.BucketID, path string, options *crud.Options) ([]byte, error) {
	ep := filesOp.endpoints[files.EPRead]
	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path

	var data []byte
	return data, server_http.Request(serverURL, ep, nil, &data, options.GetIdentity(), filesOp.logfile)
}

func (filesOp *filesHTTP) Remove(bucketID files.BucketID, path string, options *crud.Options) error {
	ep := filesOp.endpoints[files.EPRemove]
	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path

	return server_http.Request(serverURL, ep, nil, nil, options.GetIdentity(), filesOp.logfile)
}

func (filesOp *filesHTTP) List(bucketID files.BucketID, path string, depth int, options *crud.Options) (files.FilesInfo, error) {
	if path = strings.TrimSpace(path); path == "" {
		path = "."
	}

	ep := filesOp.endpoints[files.EPList]
	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + strconv.Itoa(depth)

	var filesInfo files.FilesInfo
	return filesInfo, server_http.Request(serverURL, ep, nil, &filesInfo, options.GetIdentity(), filesOp.logfile)
}

func (filesOp *filesHTTP) Stat(bucketID files.BucketID, path string, depth int, options *crud.Options) (*files.FileInfo, error) {
	if path = strings.TrimSpace(path); path == "" {
		path = "."
	}

	ep := filesOp.endpoints[files.EPStat]
	serverURL := filesOp.host + ep.Path + "/" + string(bucketID) + "/" + path + "/" + strconv.Itoa(depth)

	var fileInfo *files.FileInfo
	return fileInfo, server_http.Request(serverURL, ep, nil, &fileInfo, options.GetIdentity(), filesOp.logfile)
}

//// logger.OperatorComments -----------------------------------------------------------
//
//var _ logger.OperatorComments = &filesHTTP{}
//
//func (filesOp *filesHTTP) Comment(text string) {
//	logger.LogIntoFile(filesOp.logfile, l, text, "on filesHTTP.Comment()")
//}
//
