package views_html_old

import "html"

func AttributesHTML(attributes Attributes) string {
	var attributesHTML string
	for k, v := range attributes {
		attributesHTML += " " + html.EscapeString(k) + `="` + html.EscapeString(v) + `"`
	}

	return attributesHTML
}

func AttributeHTML(key, value string) string {
	return key + `="` + html.EscapeString(value) + `"`
}
