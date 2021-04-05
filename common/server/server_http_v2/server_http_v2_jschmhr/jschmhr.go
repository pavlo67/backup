package server_http_v2_jschmhr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server/server_http"

	server_http_v2 "github.com/pavlo67/tools/common/server/server_http_v2"
)

var _ server_http_v2.OperatorV2 = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	port        int
	tlsCertFile string
	tlsKeyFile  string

	// onRequest server_http_v2.OnRequestMiddleware

	wrappersHTTP map[server_http_v2.WrapperHTTPKey]server_http_v2.WrapperHTTP
}

func New(port int, tlsCertFile, tlsKeyFile string, onRequest server_http.OnRequestMiddleware, wrappersHTTP map[server_http_v2.WrapperHTTPKey]server_http_v2.WrapperHTTP) (server_http_v2.OperatorV2, error) {
	if port <= 0 {
		return nil, fmt.Errorf("on server_http_jschmhr.New(): wrong port = %d", port)
	}

	//if onRequest == nil {
	//	return nil, errors.New("on server_http_jschmhr.New(): no server_http_v2.OnRequestMiddleware")
	//}

	router := httprouter.New()

	return &serverHTTPJschmhr{
		httpServer: &http.Server{
			Handler:        router,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		httpServeMux: router,
		port:         port,
		tlsCertFile:  tlsCertFile,
		tlsKeyFile:   tlsKeyFile,

		// onRequest: onRequest,
		wrappersHTTP: wrappersHTTP,
	}, nil
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *serverHTTPJschmhr) Start() error {
	if s == nil {
		return errors.New("no serverOp to start")
	}

	s.httpServer.Addr = ":" + strconv.Itoa(s.port)
	l.Info("Server is starting on address ", s.httpServer.Addr)

	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		return s.httpServer.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
	}

	return s.httpServer.ListenAndServe()
}

func (s *serverHTTPJschmhr) Addr() (port int, https bool) {
	return s.port, s.tlsCertFile != "" && s.tlsKeyFile != ""
}

const onHandle = "on serverHTTPJschmhr.Handle()"

func (s *serverHTTPJschmhr) Handle(key server_http.EndpointKey, serverPath string, wrapperHTTPKey server_http_v2.WrapperHTTPKey, data interface{}) error {
	wrapperHTTP := s.wrappersHTTP[wrapperHTTPKey]
	if wrapperHTTP == nil {
		return fmt.Errorf(onHandle+": wrong wrapperHTTPKey (%s) on %s [%s]", wrapperHTTPKey, key, serverPath)
	}

	method, path, handler, err := wrapperHTTP(s, serverPath, data)
	if err != nil {
		return err
	}

	if handler == nil {
		return errors.New(onHandle + ": " + method + ": " + path + "\t!!! NULL workerHTTP ISN'T DISPATCHED !!!")
	}

	s.HandleOptions(key, path)

	l.Infof("%-10s: %s %s", key, method, path)
	switch method {
	case "GET":
		s.httpServeMux.GET(path, handler)
	case "POST":
		s.httpServeMux.POST(path, handler)
	case "PUT":
		s.httpServeMux.PUT(path, handler)
	case "DELETE":
		s.httpServeMux.DELETE(path, handler)
	default:
		return fmt.Errorf(onHandle+": method (%s) isn't supported", method)
	}

	return nil
}

func (s *serverHTTPJschmhr) HandleOptions(key server_http.EndpointKey, serverPath string) {
	s.httpServeMux.OPTIONS(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l.Infof("%-10s: OPTIONS %s", key, serverPath)
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
	})
}

var reHTMLExt = regexp.MustCompile(`\.html?$`)

const filepathSuffix = "*filepath"

func (s *serverHTTPJschmhr) HandleFiles(key server_http.EndpointKey, serverPath string, staticPath server_http.StaticPath) error {
	l.Infof("%-10s: FILES %s <-- %s", key, serverPath, staticPath.LocalPath)

	// TODO: check localPath

	if staticPath.MIMEType == nil {

		// TODO!!! CORS
		if len(serverPath) < len(filepathSuffix) || serverPath[len(serverPath)-len(filepathSuffix):] != filepathSuffix {
			serverPath += "*filepath"
		}

		s.httpServeMux.ServeFiles(serverPath, http.Dir(staticPath.LocalPath))
		return nil
	}

	s.HandleOptions(key, serverPath)

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		if staticPath.MIMEType != nil && *staticPath.MIMEType != "" {
			w.Header().Set("Content-Type", *staticPath.MIMEType)
		}

		OpenFile, err := os.Open(staticPath.LocalPath + "/" + p.ByName("filepath"))
		defer OpenFile.Close()
		if err != nil {
			l.Error(err)
		} else {
			io.Copy(w, OpenFile)
		}

		//if mimeType != nil {
		//}
		//fileServer.ServeHTTP(w, r)
	})

	return nil
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.ExportID.PathWithParams, nil)
// if err != nil {
//	l.ErrStr("can't read MIMEType for file: ", localPath+"/"+r.ExportID.PathWithParams, err)
// }
