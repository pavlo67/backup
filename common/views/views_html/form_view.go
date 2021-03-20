package views_html

import "html"

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

// view - not editable data field
// text - text label only (no data field linked to!)

func FieldView(field Field, data map[string]string) (string, string) { // , frontOps map[string]Operator

	//if frontOp, ok := frontOps[field.Type]; ok {
	//	params := map[string]string{
	//		"format": field.Format,
	//		"class":  field.Class,
	//		"style":  "width:100%",
	//	}
	//	return html.EscapeString(field.Label), frontOp.HTMLToView(field, data[field.Key], nil, params)
	//}

	var types = []string{"", "view", "text", "select", "checkbox", "date"}
	isType := false
	for _, v := range types {
		if v == field.Type {
			isType = true
			break
		}
	}
	if !isType {
		return "", ""
	}

	var class = field.Attributes["class"]
	if class != "" {
		class = ` class="` + html.EscapeString(class) + `"`
	}

	var resHTML string

	if field.Options.StringDefault("format", "") == "datetime" {
		resHTML = html.EscapeString(data[field.Key])

	} else if field.Options.StringDefault("format", "") == "url" {
		var url = html.EscapeString(data[field.Key])
		resHTML = `<a href="` + url + `" target=_blank>` + url + `</a>`

	} else if field.Type == "text" {
		resHTML = html.EscapeString(field.Options.StringDefault("format", ""))

	} else if field.Type == "checkbox" {
		if data[field.Key] == "on" {
			resHTML = "так"
		} else if field.Attributes["class"] != "not_empty" {
			resHTML = "ні"
		}
	} else if field.Attributes["class"] == "not_empty" && data[field.Key] == "0" {
		// shows nothing
	} else {
		resHTML = html.EscapeString(data[field.Key])

	}
	return html.EscapeString(field.Label), resHTML
}
