package notebook_html

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data_exchange/components/ns"
	"github.com/pavlo67/data_exchange/components/tags"
	server_http "github.com/pavlo67/tools/common/server/server_http2"
	"github.com/pavlo67/tools/entities/records"
)

const No = `<b style="size:18px;">»</b> `

//const My = `<b style="color:rgb(255,0,127);size:18px;">»</b> `
//const Group = `<b style="color:rgb(204,204,0);size:18px;">»</b> `
//const Public = `<b style="color:rgb(63,255,63);size:18px;">»</b> `

const MyImgPath = "/images/my.jpg"
const GroupImgPath = "/images/group.jpg"
const PublicImgPath = "/images/free.jpg"

const My = `<img src="/images/my.jpg" title="приватні сторінки"> `
const Group = `<img src="/images/group.jpg" title="сторінки групи"> `
const Public = `<img src="/images/free.jpg" title="сторінки для загалу"> `

const Tags = `<img src="/images/tags.png" title="рубрикатор" style="vertical-align:-5%;">`

const ImgFile = "/images/file.png"
const ImgPlus = "/images/plus.png"
const ImgMinus = "/images/minus.png"

const listDelimiter = "<hr>"

const txtNoTags = "Нема міток"

func HTMLAuthor(r *records.Item, identity *auth.Identity) string {
	info := "створено: " + r.CreatedAt.Format("02.01.2006 15:04")
	if r.UpdatedAt != nil {
		info += `, оновлено: ` + r.UpdatedAt.Format("02.01.2006 15:04:05")
	}

	//var image string
	//if r.RView == r.ROwner {
	//	image = `<img src="` + _views.MyImgPath + `" title="` + info + `">`
	//} else if r.RView == controller.Anyone {
	//	image = `<img src="` + _views.PublicImgPath + `" title="` + info + `">`
	//} else {
	//	image = `<img src="` + _views.GroupImgPath + `" title="` + info + `">`
	//}
	//if r.ROwner == userIdentity.String() {
	//	image = `<a href="` + endpoints["edit"].Path(r.ID) + `">` + image + `</a>`
	//}
	//var author string
	//if r.Author == "" {
	//	author = `<a href="` + endpoints["items"].Path(string(r.ROwner)) + `">???</a>`
	//} else {
	//	author = `<a href="` + endpoints["items"].Path(string(r.ROwner)) + `">` + r.Author + `</a>`
	//}
	//
	//return image + "<small>Автор: <b>" + author + "</b><br>" + r.Visibility + "</small>\n"

	return info
}

//func HTMLParents(parents []records.Item) string {
//	if len(parents) < 1 {
//		return ""
//	}
//
//	if len(parents) == 1 {
//		return `В розділі: <a href="` + endpoints["view"].Path(parents[0].ID) + `">` + parents[0].Name + `</a>`
//	}
//
//	var htmlLinks string
//	for _, p := range parents {
//		htmlLinks += `<br>- <a href="` + endpoints["view"].Path(p.ID) + `">` + p.Name + `</a>`
//	}
//
//	return "В розділах:" + htmlLinks
//}

func HTMLTags(tags []tags.Item, ViewerNSS, OwnerNSS ns.NSS, epTagged server_http.Get1, joiner string) string {
	var htmlTags string

	for _, tag := range tags {
		//objectID := strings.TrimSpace(l.To)
		//if objectID != "" && objectID != "0" {
		//	link = endpoints["view"].Path(objectID)
		//} else if OwnerNSS != "" {
		//	link = endpoints["itemsByTagOwn"].Path(string(OwnerNSS), tag)
		//} else {
		//	link = endpoints["itemsByTag"].Path(string(ViewerNSS), tag)
		//}

		urlStr, err := epTagged(tag)
		if err != nil || urlStr == "" {
			l.Errorf("can't htmlOp.epTagged(%s), got %s, %s", tag, urlStr, err)
		}
		htmlTags += joiner + `<a href="` + urlStr + `">` + tag + "</a>\n"
	}

	if htmlTags == "" {
		return ""
	}

	return "<small>В розділах:</small>" + htmlTags
}

//func HTMLTagsList(tags []links.TagInfo, rView, rOwner *auth.Identity, userIsAdmin bool) string {
//	if len(tags) < 1 {
//		return txtNoTags
//	}
//
//	var link, htmlTags string
//
//	for _, t := range tags {
//
//		setParent := ""
//		if userIsAdmin {
//			//setParent = "<a href=\"" + nt.endpoints["editSection"].ServerPath + urlTag + "\"><span style=\"color:orange;\"> ✎ </span></a>"
//		}
//
//		tag := html.EscapeString(t.Tag)
//
//		if rOwner != "" {
//			link = endpoints["itemsByTagOwn"].Path(string(rOwner), tag)
//		} else {
//			link = endpoints["itemsByTag"].Path(string(rView), tag)
//		}
//
//		htmlTags += `<li>` + setParent +
//			`<a href="` + link + `">` + tag + "</a>  " +
//			"[" + strconv.FormatUint(t.Count, 10) + "]" +
//			"</li>\n"
//	}
//
//	return htmlTags
//}
//
//func HTMLTagsIndex(userIdentity *confidenter.Identity, rView, rOwner *auth.Identity) string {
//	if rOwner != "" {
//		rOwnerNSSentity := rOwner.Identity()
//		if userIdentity != nil && rOwnerNSSentity.ID == userIdentity.ID {
//			return `<a href="` + endpoints["tagsOwn"].Path(string(rOwner)) + `">` + _views.Tags + " всі мої мітки</a>\n\n<p>"
//		} else {
//			return `<a href="` + endpoints["tagsOwn"].Path(string(rOwner)) + `">` + _views.Tags + " всі мітки користувача</a>\n\n<p>"
//		}
//	}
//
//	rViewIdentity := rView.Identity()
//	if rViewIdentity.Path == "user" {
//		if userIdentity != nil && rViewIdentity.ID == userIdentity.ID {
//			return `<a href="` + endpoints["tagsOwn"].Path(string(rView)) + `">` + _views.Tags + " всі мої мітки</a>\n\n<p>"
//		}
//		return `<a href="` + endpoints["tagsOwn"].Path(string(rView)) + `">` + _views.Tags + " всі мітки користувача</a>\n\n<p>"
//	} else if rViewIdentity.Path == "group" {
//		return `<a href="` + endpoints["tags"].Path(string(rView)) + `">` + _views.Tags + " всі мітки групи</a>\n\n<p>"
//	}
//
//	return `<a href="` + endpoints["tags"].Path(string(controller.Anyone)) + `">` + _views.Tags + " всі мітки для загалу</a>\n\n<p>"
//}
