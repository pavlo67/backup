package notebook_vews_old

//// share ------------------------------------------------------------------------------------------------------
//
//const txtNoRecords = "Нема записів"
//const publicChangesPostfix = "/change"
//
//var ReCommonEdit = regexp.MustCompile(publicChangesPostfix + `$`)
//
//func dataView(user *auth.User, rView auth.ID, publicChanges bool) ([][2]string, string) {
//	values := [][2]string{
//		{things_old.Private, string(controller.Owner)},
//		{things_old.Public, string(basis.Anyone)},
//	}
//	for _, t := range user.Accesses {
//		if t.Right == rights.Member {
//			values = append(values, [2]string{t.Label, string(t.IS)})
//			values = append(values, [2]string{t.Label + " (спільне редаґування)", string(t.IS) + publicChangesPostfix})
//		}
//	}
//	selected := string(rView)
//	if publicChanges {
//		selected += publicChangesPostfix
//	}
//	return values, selected
//}
//
//func View(fields []viewshtml.Field, data map[string]string, options map[string]viewshtml.SelectString, frontOps map[string]viewshtml.Operator) string {
//	if data == nil {
//		data = map[string]string{}
//	}
//	var viewHTML, titleHTML, resHTML string
//	for _, f := range fields {
//		titleHTML, resHTML = viewshtml.FieldView(f, data, options, frontOps)
//		if resHTML == "" && f.Params[viewshtml.NotEmptyKey] == true {
//			continue
//		}
//
//		if titleHTML != "" {
//			titleHTML = "<small>" + titleHTML + ":</small>\n"
//		}
//		viewHTML += "<div>\n" + titleHTML + resHTML + "\n</div>\n"
//	}
//
//	return viewHTML
//	// +`<input id=links_list type=hidden value="` + html.EscapeString(data["tags"]) + `">` + "\n"
//}
//
////htmlEdit := "\n<form enctype=\"multipart/form-data\">\n"
////</form>
//
//var formNum = 0
//
//const formNameBase = "editst_"
//
//func Edit(user *auth.User, fields []viewshtml.Field, data map[string]string, options map[string]viewshtml.SelectString, frontOps map[string]viewshtml.Operator, rView auth.ID, publicChanges bool) string {
//	formNum++
//	formID := formNameBase + strconv.Itoa(formNum) + "_"
//
//	if data == nil {
//		data = map[string]string{}
//	}
//	if options == nil {
//		options = map[string]viewshtml.SelectString{}
//	}
//
//	options["visibility"], data["visibility"] = dataView(user, rView, publicChanges)
//
//	var editHTML, titleHTML, resHTML string
//	for _, f := range fields {
//		titleHTML, resHTML = viewshtml.FieldEdit(formID, f, data, options, frontOps)
//
//		if resHTML == "" && f.Params[viewshtml.NotEmptyKey] == true {
//			continue
//		}
//
//		if titleHTML != "" {
//			titleHTML = "<small>" + titleHTML + ":</small> \n"
//		}
//
//		if f.Params[viewshtml.InlineFields] != true {
//			resHTML += `<div class="ut">` + "\n</div>"
//		}
//
//		editHTML += `<div id="div_` + html.EscapeString(formID+f.Key) + `">` + "\n" +
//			titleHTML + resHTML +
//			"</div>\n"
//
//	}
//
//	return editHTML
//}
//
//func HTMLItemIndex(genusKey, id string, rView auth.ID, blank, view, edit, remove bool) string {
//	var htmlItemIndex string
//	if blank {
//		if rView == "" {
//			rView = basis.Anyone
//		}
//		htmlItemIndex += `- <a href="` + endpoints["blank"].Path(genusKey, string(rView)) + `">новий запис</a>` +
//			"</td></tr>\n<tr><td>"
//	}
//
//	if view && id != "" {
//		htmlItemIndex += `- <a href="` + endpoints["view"].Path(id) + `">перегляд без збереження змін</a>` + "<br>\n"
//	}
//
//	if edit && id != "" {
//		htmlItemIndex += `- <a href="` + endpoints["edit"].Path(id) + `">редаґування</a>` + "<br>\n"
//	}
//
//	if remove && id != "" {
//		htmlItemIndex += `- <a href="#" id="` + listeners["delete"].ID + `" data-id="` + html.EscapeString(id) + `">вилучити запис</a>` + "<br>\n"
//	}
//
//	if htmlItemIndex == "" {
//		return ""
//	}
//
//	return "<tr><td>" + htmlItemIndex + "</td></tr>\n"
//}
//
//const LoadableID = "loadable_"
//
//var num int
//
//func HTMLLoadable(title, topID, genus, objectID string, dataDefault genera.DataRaw, visible bool) string {
//	var htmlContent string
//
//	num++
//
//	idNum := "_" + strconv.Itoa(num)
//	if objectID == "" {
//		objectID = "blank"
//	}
//	general := ` data-top_id="` + topID + `" data-genus="` + genus + `" data-id="` + objectID + `"`
//
//	if len(dataDefault) > 0 {
//		objects := map[string]genera.DataRaw{
//			"default": dataDefault,
//		}
//		objectsStr, err := json.Marshal(objects)
//		if err != nil {
//			log.Printf("on json.Marshal(%#v): %s", objects, err)
//		} else {
//			general += ` data-objects="` + html.EscapeString(string(objectsStr)) + `"`
//		}
//	}
//
//	// log.Print(general)
//
//	if visible {
//		htmlContent = `<img class="clicked" id="` + LoadableID + "img" + idNum + `" src="` + viewshtml.ImgMinus + `"` + general + `> ` +
//			`<span id="` + LoadableID + "title" + idNum + `">` + title + "</span>" +
//			`<p id="` + LoadableID + "content" + idNum + `"></p>`
//		// add JS to "click" on the image (it will show some content)
//
//	} else {
//		htmlContent = `<img class="clicked" id="` + LoadableID + "img" + idNum + `" src="` + viewshtml.ImgPlus + `"` + general + `> ` +
//			`<span id="` + LoadableID + "title" + idNum + `">` + title + "</span>" +
//			`<p id="` + LoadableID + "content" + idNum + `" style="visibility:hidden;"></p>`
//
//	}
//
//	return htmlContent
//}

// titleID := "loadable_title_" + genus + "_" + objectID
// contentID := "loadable_content_" + genus + "_" + objectID
// 				`<span objectID="` + titleID + `">` + title + "</span>\n" +
//				`<br><div objectID="` + contentID + `" style="visibility:hidden;position:absolute;"></div>`
