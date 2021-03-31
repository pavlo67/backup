//// +build fullserver linux
//
package notebook_vews_old

//
//import (
//	"net/http"
//
//	"github.com/pavlo67/punctum/interfaces/confidenter"
//	"github.com/pavlo67/punctum/interfaces/controller"
//	"github.com/pavlo67/punctum/interfaces/controller/rights"
//
//	"github.com/pavlo67/punctum/components/_views"
//)
//
//var tplLeftNoUser map[string]string
//var htmlLeftTop, htmlLeftBottom string
//
//func initHTML() {
//	htmlPublic := `<div class="ul">` + _views.Public +
//		`<a href="` + itemsEndpoints["tags"].Path(string(controller.Anyone)) + `">` + _views.Tags + `</a> ` +
//		`<a href="` + itemsEndpoints["itemsTop"].Path() + `">записи для загалу</a>` + "</div>\n"
//
//	htmlSearch := `<div class="ul utd">` +
//		`<input style="width:171px;" id="to_search">` +
//		`<input type="button" id="` + listeners["search"].ID + `" value="знайти">` +
//		"</div>\n\n"
//
//	htmlSearchMy := `<div class="ul utd">` +
//		`<input style="width:171px;" id="to_search_my"">` +
//		`<input type="button" id="` + listeners["searchMy"].ID + `" value="знайти">` +
//		"</div>\n\n"
//
//	htmlLeftTop = `<div class="ut">` + _views.Public + `<a href="/">Ой, мамо! Де я???</a></div>` +
//		`<div class="ut">` + _views.My + ` Мої сторінки</div>` +
//		`<div class="ul">` + _views.My +
//		`<a href="` + itemsEndpoints["tagsMy"].ServerPath + `">` + _views.Tags + `</a> ` +
//		`<a href="` + itemsEndpoints["itemsMy"].ServerPath + `">записи</a>` + "</div>\n" +
//		`<div class="ul">` + _views.My + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(controller.Owner)) + `">новий запис</a></div>` +
//		htmlSearchMy +
//		`<div class="ut">` + _views.Public + ` Публічні сторінки` +
//		htmlPublic
//
//	htmlLeftBottom = "\n</div>" +
//		`<div class="ul">` + _views.Public + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(controller.Anyone)) + `">новий запис</a></div>` +
//		"\n</div>" +
//		htmlSearch
//
//	htmlLeftNoUser := `<div class="ut gray">` + _views.No + ` Мої сторінки</div>` +
//		`<div class="ul gray">` + _views.No + _views.Tags + "записи</a></div>\n" +
//		`<div class="ul gray">` + _views.No + ` новий запис</a></div>` +
//		`<div class="ut">` + _views.Public + ` Публічні сторінки` +
//		htmlPublic +
//		htmlSearch
//
//	tplLeftNoUser = map[string]string{
//		"left":  htmlLeftNoUser,
//		"front": htmlFront,
//	}
//}
//
//func notebookTemplator(r *http.Request, user *confidenter.User) map[string]string {
//	if user == nil || user.ID == "" {
//		return tplLeftNoUser
//	}
//
//	htmlLeft := htmlLeftTop
//
//	for _, gr := range user.Accesses {
//		if gr.Right == rights.Member {
//			htmlLeft += `<div class="ul">` + _views.Group +
//				`<a href="` + itemsEndpoints["tags"].Path(string(gr.IdentityString)) + `">` + _views.Tags + `</a> ` +
//				`<a href="` + itemsEndpoints["items"].Path(string(gr.IdentityString)) + `">` + gr.Label + `</a></div>`
//		}
//	}
//
//	tplLeft := map[string]string{
//		"left":  htmlLeft + htmlLeftBottom,
//		"front": htmlFront,
//	}
//
//	return tplLeft
//}
