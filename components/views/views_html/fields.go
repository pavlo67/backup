package views_html

import (
	"html"
	"regexp"
	"strings"
)

type Field struct {
	Key    string
	Label  string
	Type   string
	Format string
	Class  string
}

var ReDigitsOnly = regexp.MustCompile(`^\d+$`)

type SelectString [][2]string

func HTMLSelect(idDOMEscaped, class string, values SelectString, selected string) string {
	body := ""
	var option string
	for i := 0; i < len(values); i++ {
		body += "<option"
		if values[i][1] != "" {
			option = values[i][1]
			body += ` value="` + html.EscapeString(values[i][1]) + `"`
		} else {
			option = values[i][0]
		}
		if option == selected {
			body += " selected"
		}
		body += ">" + html.EscapeString(values[i][0]) + "</option>\n"
	}
	return `<select id="` + idDOMEscaped + `" ` + class + `>` + body + "</select>\n"
}

func GeneralEditParts(formID, idDOM, class string) (string, string, string, string, string) {
	idDOMNoFormEscaped := html.EscapeString(idDOM)
	if formID != "" {
		idDOM = formID + "_" + idDOM
	}

	if class != "" {
		class = ` class="` + html.EscapeString(class) + `"`
	}

	idDOMEscaped := html.EscapeString(idDOM)
	var general = `id="` + idDOMEscaped + `"` + class

	var generalNoForm = `id="` + idDOMNoFormEscaped + `"` + class
	// to add listener by element.Id

	return idDOM, idDOMEscaped, class, general, generalNoForm

}

func FieldEdit(formID string, field Field, data map[string]string, values map[string]SelectString) (string, string) { // , frontOps map[string]Operator

	if field.Type == "view" {
		return FieldView(field, data) // , frontOps
	}

	_, idDOMEscaped, class, general, generalNoForm := GeneralEditParts(formID, field.Key, field.Class)

	var titleHTML = html.EscapeString(field.Label)
	var resHTML string

	if field.Type == "password" {
		resHTML = `<input style="width:100%" type="password" ` + general + ` />`
	} else if field.Type == "select" {
		resHTML = HTMLSelect(idDOMEscaped, class, values[field.Key], data[field.Key])
	} else if field.Type == "text" {
		resHTML = html.EscapeString(field.Format)
	} else if field.Type == "checkbox" {
		var checked string
		if data[field.Key] != "" {
			checked = " checked"
		}
		resHTML = `<input type="checkbox" ` + general + checked + `/>`
		//} else if frontOp, ok := frontOps[field.Type]; ok {
		//	params := map[string]string{
		//		"form_id": formID,
		//		"style":   "width:100%",
		//	}
		//	resHTML = frontOp.HTMLToEdit(field, data[field.Key], values[field.Key], params)
	} else {
		var value = html.EscapeString(data[field.Key])
		if field.Type == "button" {
			resHTML = `<input type="button" ` + generalNoForm + ` data-form_id="` + html.EscapeString(formID) + `" data-value="` + value + `" value="` + titleHTML + `" />`
			titleHTML = ""
		} else if field.Type == "hidden" {
			resHTML = `<input type="hidden" ` + general + ` value="` + value + `" /> `
			titleHTML = ""
		} else if ReDigitsOnly.MatchString(strings.TrimSpace(field.Format)) {
			resHTML = `<textarea style="width:100%" ` + general + ` rows=` + field.Format + `>` + value + `</textarea>`
		} else if field.Type == "date" {
			resHTML = `<input type="date" ` + general + ` value="` + value + `" />`
		} else {
			resHTML = `<input type="` + field.Type + `"style="width:100%" ` + general + ` value="` + value + `" />`
		}
	}
	return titleHTML, resHTML
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

	var class = field.Class
	if class != "" {
		class = " Class=\"" + html.EscapeString(class) + "\""
	}

	var resHTML string

	if field.Format == "datetime" {
		resHTML = html.EscapeString(data[field.Key])

	} else if field.Format == "url" {
		var url = html.EscapeString(data[field.Key])
		resHTML = `<a href="` + url + `" target=_blank>` + url + `</a>`

	} else if field.Type == "text" {
		resHTML = html.EscapeString(field.Format)

	} else if field.Type == "checkbox" {
		if data[field.Key] == "on" {
			resHTML = "так"
		} else if field.Class != "not_empty" {
			resHTML = "ні"
		}
	} else if field.Class == "not_empty" && data[field.Key] == "0" {
		// shows nothing
	} else {
		resHTML = html.EscapeString(data[field.Key])

	}
	return html.EscapeString(field.Label), resHTML
}
