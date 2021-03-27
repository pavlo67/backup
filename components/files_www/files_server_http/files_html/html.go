package files_html

import (
	"github.com/pavlo67/tools/entities/files"

	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

var _ Operator = &filesHTML{}

type filesHTML struct {
	//epCreate string
	//epView   server_http.Get1
}

const onNew = "on filesHTML.New(): "

func New(pagesConfig server_http.ConfigPages) (Operator, error) { // , restConfig

	//epCreate, err := server_http.CheckGet0(pagesConfig, files_www.IntefaceKeyHTMLCreate, false)
	//if err != nil {
	//	return nil, err
	//}
	//
	//epView, err := server_http.CheckGet1(pagesConfig, files_www.IntefaceKeyHTMLView, false)
	//if err != nil {
	//	return nil, err
	//}

	return &filesHTML{
		//epCreate: epCreate,
		//epView:   epView,
	}, nil
}

// TODO!!! look at https://github.com/kataras/blocks

func (htmlOp *filesHTML) FragmentsList(filesItems []files.Item, path string, identity *auth.Identity) (server_http.Fragments, error) {
	htmlFiles := htmlOp.HTMLFiles(filesItems, identity)

	fragments := server_http.Fragments{
		"title":   "каталог: " + path,
		"header":  "Каталог: " + path,
		"content": htmlFiles,
	}

	return fragments, nil
}

func (htmlOp *filesHTML) HTMLFiles(filesItems []files.Item, identity *auth.Identity) string {
	if len(filesItems) < 1 {
		return "нема файлів"
	}

	var htmlFiles string

	for _, f := range filesItems {
		details := ""
		name := f.Path
		urlStr := ""

		//urlStr, err := htmlOp.epView(string(f.ID))
		//if err != nil || urlStr == "" {
		//	l.Errorf("can't htmlOp.epView(%s), got %s, %s", f.ID, urlStr, err)
		//}

		htmlFiles += `<li><a href="` + urlStr + `">` + name + "</a></li>\n" +
			"<br>" + details + // HTMLHidden(details) +
			"\n"
	}

	return htmlFiles
}
