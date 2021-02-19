package views_html

import (
	"html"
	"regexp"
	"strings"
)

var ReDigitsOnly = regexp.MustCompile(`^\d+$`)

func GeneralEditPart(formID, fieldKey string, attributes map[string]string) string {
	return ` id="` + formID + fieldKey + `" name="` + fieldKey + `" ` + AttributesHTML(attributes)

}

func FieldEdit(formID string, field Field, data map[string]string, values map[string]SelectString) (string, string) { // , frontOps map[string]Operator

	if field.Type == "view" {
		return FieldView(field, data) // , frontOps
	}

	general := GeneralEditPart(formID, field.Key, field.Attributes) // generalNoForm

	var titleHTML, resHTML string

	if field.Type == "password" {
		resHTML = `<input style="width:100%" type="password" ` + general + ` />`
	} else if field.Type == "select" {
		resHTML = HTMLSelect(general, values[field.Key], data[field.Key])
	} else if field.Type == "text" {
		resHTML = html.EscapeString(field.Options.StringDefault("format", ""))
	} else if field.Type == "checkbox" {
		var checked string
		if data[field.Key] != "" {
			checked = " checked"
		}
		resHTML = `<input type="checkbox" ` + checked + general + `/>`
		//} else if frontOp, ok := frontOps[field.Type]; ok {
		//	params := map[string]string{
		//		"form_id": formID,
		//		"style":   "width:100%",
		//	}
		//	resHTML = frontOp.HTMLToEdit(field, data[field.Key], values[field.Key], params)
	} else {
		var titleHTML = html.EscapeString(field.Label)
		var value = html.EscapeString(data[field.Key])
		if field.Type == "button" || field.Type == "submit" {
			resHTML = `<input type="` + field.Type + `" value="` + titleHTML + `"` + general + ` />`
			titleHTML = ""
		} else if field.Type == "hidden" {
			resHTML = `<input type="hidden"  value="` + value + `"` + general + `/>`
			titleHTML = ""
		} else if format := strings.TrimSpace(field.Options.StringDefault("format", "")); ReDigitsOnly.MatchString(format) {
			resHTML = `<textarea style="width:100%"  rows=` + format + general + `>` + value + `</textarea>`
		} else if field.Type == "date" {
			resHTML = `<input type="date" value="` + value + `"` + general + `/>`
		} else {
			var fieldType string
			if field.Type != "" {
				fieldType = ` type="` + field.Type + `"`
			}
			resHTML = `<input` + fieldType + ` style="width:100%"  value="` + value + `"` + general + ` />`
		}
	}
	return titleHTML, resHTML
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
