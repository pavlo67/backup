package tags_html

import (
	"fmt"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/server/server_http"
	server_http2 "github.com/pavlo67/tools/common/server/server_http"

	"github.com/pavlo67/tools/components/notebook"
	"github.com/pavlo67/tools/components/tags"
)

var _ Operator = &tagsHTML{}

// should be thread-safe
type tagsHTML struct {
	pagesConfig server_http.Config
	restConfig  server_http.Config
}

const onNew = "on tagsHTML.New(): "

func New(pagesConfig, restConfig server_http.Config) (Operator, error) {
	method, urlStr, err := server_http2.EP(pagesConfig, notebook.IntefaceKeyHTMLTagged, []string{"tag"}, false)
	if err != nil || urlStr == "" {
		return nil, fmt.Errorf("can't EP(%#v, notebook.IntefaceKeyHTMLTags, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	} else if strings.TrimSpace(strings.ToUpper(method)) != "GET" {
		return nil, fmt.Errorf("wrong method on EP(%#v, notebook.IntefaceKeyHTMLTags, nil, false), got %s, %s, %s", pagesConfig, method, urlStr, err)
	}

	return &tagsHTML{
		pagesConfig: pagesConfig,
		restConfig:  restConfig,
	}, nil
}

const onPrepare = "on tagsHTML.Prepare(): "

func (htmlOp *tagsHTML) Prepare(key Key, template string, params common.Map) error {
	return nil
}

const onHTMLTags = "on tagsHTML.HTMLTags(): "

func (htmlOp *tagsHTML) HTMLTags(tss tags.Stats) (string, error) {
	htmlTags := "\n\n<ul>\n"

	for i, ts := range tss {
		tag := strings.TrimSpace(string(ts.Tag))
		if tag == "" {
			l.Errorf("empty tag %d in %#v", i, tss)
			continue
		}

		method, urlStr, err := server_http2.EP(htmlOp.pagesConfig, notebook.IntefaceKeyHTMLTagged, []string{tag}, false)
		if err != nil || urlStr == "" {
			l.Errorf("can't server_http.EP(%#v, notebook.IntefaceKeyHTMLTags,nil,false), got %s, %s, %s", htmlOp.pagesConfig, method, urlStr, err)
			continue
		}

		htmlTags += fmt.Sprintf(`<li><a href="%s">%s</a> [%d]</li>`, urlStr, tag, ts.Count)
	}

	htmlTags += "</ul>\n\n"

	return htmlTags, nil
}
