package files_server_http

import (
	"path"

	"github.com/pavlo67/tools/entities/files"

	"github.com/pavlo67/common/common/auth"
	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

type filesHTML struct {
	//epCreate string
	//epView   server_http.Get1
}

const onNew = "on filesHTML.New(): "

func New(pagesConfig server_http.ConfigPages) (*filesHTML, error) { // , restConfig

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

func (htmlOp *filesHTML) FragmentsList(prefix, basePath string, filesItems []files.Item, path string, identity *auth.Identity) (server_http.Fragments, error) {
	htmlFiles := htmlOp.HTMLFiles(prefix, basePath, filesItems, identity)

	fragments := server_http.Fragments{
		"title":   "каталог: " + path,
		"header":  "Каталог: " + path,
		"content": htmlFiles,
	}

	return fragments, nil
}

func (htmlOp *filesHTML) HTMLFiles(prefix, basePath string, filesItems []files.Item, identity *auth.Identity) string {
	var htmlFiles string

	if dir := path.Dir(basePath); dir != "." && dir != basePath {
		urlStr := "/" + prefix + "/list" + dir
		htmlFiles += `<li><a href="` + urlStr + `">..` + "</a></li>\n"

	}

	//if len(filesItems) < 1 {
	//	return "нема файлів"
	//}

	for _, f := range filesItems {
		details := ""
		name := f.Path
		urlStr := ""
		if f.IsDir {

			// TODO!!! use PagesConfig
			urlStr = "/" + prefix + "/list" + basePath + f.Path
		}

		//urlStr, err := htmlOp.epView(string(f.ID))
		//if err != nil || urlStr == "" {
		//	l.Errorf("can't htmlOp.epView(%s), got %s, %s", f.ID, urlStr, err)
		//}

		htmlFiles += `<li><a href="` + urlStr + `">` + name + "</a></li>\n" +
			details + // HTMLHidden(details) +
			"\n"
	}

	return htmlFiles
}
