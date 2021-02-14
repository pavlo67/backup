package records_html

import (
	"fmt"
	"strings"

	"github.com/pavlo67/tools/components/tags"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/server/server_http"
	server_http2 "github.com/pavlo67/tools/common/server/server_http"

	"github.com/pavlo67/tools/components/formatter"
	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/records"
)

var _ Operator = &recordsHTML{}

// should be thread-safe
type recordsHTML struct {
	pagesConfig server_http.Config
	restConfig  server_http.Config
}

const onNew = "on recordsHTML.New(): "

func New(pagesConfig, restConfig server_http.Config) (Operator, error) {
	method, urlStr, err := server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLView, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLView, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLTagged, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLTagged, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLTagged, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	method, urlStr, err = server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLEdit, []string{"id"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLEdit, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLEdit, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	return &recordsHTML{
		pagesConfig: pagesConfig,
		restConfig:  restConfig,
	}, nil
}

const onPrepare = "on recordsHTML.Prepare(): "

func (formatterOp *recordsHTML) Prepare(key formatter.Key, template string, params common.Map) error {
	return nil
}

const onHTMLView = "on recordsHTML.HTMLView(): "

func (formatterOp *recordsHTML) HTMLView(r *records.Item, children []records.Item) (string, error) {
	return fmt.Sprintf("%#v / %#v", r, children), nil
	//htmlTags := "\n\n<ul>\n"
	//
	//for i, ts := range tss {
	//	tag := strings.TrimSpace(string(ts.Tag))
	//	if tag == "" {
	//		l.Errorf("empty tag %d in %#v", i, tss)
	//		continue
	//	}
	//
	//	method, urlStr, err := server_http2.EP(formatterOp.pagesConfig, notebook.IntefaceKeyHTMLTagged, []string{tag}, false)
	//	if err != nil || urlStr == "" {
	//		l.Errorf("can't server_http.EP(%#v, notebook.IntefaceKeyHTMLTags,nil,false), got %s, %s, %s", formatterOp.pagesConfig, method, urlStr, err)
	//		continue
	//	}
	//
	//	htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	//}
	//
	//htmlTags += "</ul>\n\n"
	//
	//return htmlTags, nil
}

func (formatterOp *recordsHTML) HTMLTagged(tag tags.Item, tagged []records.Item) (string, error) {
	return fmt.Sprintf("%s / %#v", tag, tagged), nil

}
