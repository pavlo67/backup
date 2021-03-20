// +build local

package notebook_vews_old

//
//import (
//	"net/http"
//
//	"github.com/pavlo67/punctum/components/_views"
//	"github.com/pavlo67/punctum/interfaces/confidenter"
//	"github.com/pavlo67/punctum/interfaces/controller"
//)
//
//var tplLeftNoUser, tplLeft map[string]string
//
//func initHTML() {
//	htmlSearchMy := `<div class="ul utd">` +
//		`<input style="width:171px;" id="to_search_my"">` +
//		`<input type="button" id="` + listeners["searchMy"].ID + `" value="знайти">` +
//		"</div>\n\n"
//
//	htmlLeft := `<div class="ut">` + _views.My + ` Мої сторінки</div>` +
//		`<div class="ul">` + _views.My +
//		`<a href="` + itemsEndpoints["tagsMy"].ServerPath + `">` + _views.Tags + `</a> ` +
//		`<a href="` + itemsEndpoints["itemsMy"].ServerPath + `">записи</a>` + "</div>\n" +
//
//		`<div class="ul">` + _views.My + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(controller.Owner)) + `">новий запис</a></div>` +
//		htmlSearchMy
//
//	htmlLeftNoUser := `<div class="ut gray">` + _views.No + ` Мої сторінки</div>` +
//		`<div class="ul gray">` + _views.No + _views.Tags + "записи</a></div>\n" +
//		`<div class="ul gray">` + _views.No + ` новий запис</a></div>`
//
//	tplLeftNoUser = map[string]string{
//		"left":  htmlLeftNoUser,
//		"front": htmlFront,
//	}
//
//	tplLeft = map[string]string{
//		"left":  htmlLeft,
//		"front": htmlFront,
//	}
//
//}
//
//func notebookTemplator(r *http.Request, user *confidenter.User) map[string]string {
//	if user == nil || user.ID == "" {
//		return tplLeftNoUser
//	}
//
//	return tplLeft
//}
