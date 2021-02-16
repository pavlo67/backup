package views_html

func HTMLViewTable(fields []Field, data map[string]string, options map[string]SelectString) string { // , frontOps map[string]Operator
	if data == nil {
		data = map[string]string{}
	}
	var viewHTML, titleHTML, resHTML string
	for _, f := range fields {
		titleHTML, resHTML = FieldView(f, data) // , options, frontOps

		//if resHTML == "" && ((f.Params[NotEmptyKey] == true) || (titleHTML == "")) {
		//	continue
		//}

		if titleHTML != "" {
			titleHTML = "<small>" + titleHTML + ":</small>\n"
		}
		viewHTML += "<tr><td>\n" + titleHTML + "</td><td>&nbsp;</td><td>" + resHTML + "\n</td></tr>\n"
	}

	return `<table cellspacing=0 style="padding-top:5px;">` +
		viewHTML + "</table>"
	// +`<input id=links_list type=hidden value="` + html.EscapeString(data["tags"]) + `">` + "\n"
}

func HTMLEditTable(fields []Field, formID, url string, data map[string]string, values map[string]SelectString) string { // ,
	// frontOps map[string]Operator, rView auth.ID, publicChanges bool
	//if data == nil {
	//	data = map[string]string{}
	//}
	//if values == nil {
	//	values = map[string]SelectString{}
	//}
	//
	//values["visibility"], data["visibility"] = dataView(user, rView, publicChanges)

	var editHTML, titleHTML, resHTML string
	for _, f := range fields {
		titleHTML, resHTML = FieldEdit(formID, f, data, values) // , frontOps

		//if resHTML == "" && f.Params[NotEmptyKey] == true {
		//	continue
		//}

		if titleHTML != "" {
			titleHTML = "<small>" + titleHTML + ":</small> \n"
		}
		editHTML += `<tr><td>` + titleHTML + "</td><td>" + resHTML + "</td></tr>\n"
		//  id="div_` + html.EscapeString(formID+f.Key) + `"

	}

	return `<table width="100%"><form action="` + url + `" method="POST">` + "\n" + editHTML + "\n</form>\n</table>\n"
}
