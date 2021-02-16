package views_html

import (
	"html"
	"regexp"
	"strings"

	"github.com/pavlo67/common/common"
)

type Field struct {
	Key        string
	Label      string
	Type       string
	Options    common.Map
	Attributes map[string]string
}

var ReDigitsOnly = regexp.MustCompile(`^\d+$`)

type SelectString [][2]string

func HTMLSelect(general string, values SelectString, selected string) string {
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
	return `<select ` + general + `>` + body + "</select>\n"
}

func GeneralEditPart(formID, fieldKey string, attributes map[string]string) string {
	// idDOMEscaped := formID + html.EscapeString(fieldKey)
	idDOM := formID + fieldKey

	attributesHTML := " "
	for k, v := range attributes {
		attributesHTML += html.EscapeString(k) + `="` + html.EscapeString(v) + `" `
	}

	return ` id="` + idDOM + `" name="` + fieldKey + `"` + attributesHTML

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
