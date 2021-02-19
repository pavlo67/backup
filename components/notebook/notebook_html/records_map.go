package notebook_html

import (
	"strings"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
	"github.com/pavlo67/tools/components/views/views_html"
)

var dataFields = []views_html.Field{
	{"id", "", "hidden", nil, nil},
	{"issued_id", "", "hidden", nil, nil},
	{"data_type", "", "hidden", nil, nil},
	// {"visibility", "тип", "select", nil, "ut"},
	// {"history_key", "", "hidden", nil, nil},

	{"title", "заголовок", "", nil, nil},
	{"summary", "коротко про", "", nil, nil},
	{"content_data", "опис", "", common.Map{"format": "35"}, nil},
	{"tags", "теми, розділи", "tag-it", nil, nil},
	// {"data_subtype", "", "hidden", "", ""},
	// {"embedded", "", "hidden", "", ""},
	// {"files", "завантажити файл", "file", "", "ut"},

	{"created_at", "створено", "view", common.Map{"format": "datetime"}, map[string]string{"class": "not_empty"}},
	{"updated_at", "востаннє відредаґовано", "view", common.Map{"format": "datetime"}, map[string]string{"class": "not_empty"}},
}

var createFields = append(dataFields, views_html.Field{"create", "зберегти запис", "button", nil, map[string]string{"class": "ut"}})

func value(data map[string][]string, key string) string {
	v := data[key]
	if len(v) == 1 {
		return v[0]
	} else if len(v) > 1 {
		return strings.Join(v, " ")
	}

	return ""
}

func RecordFromData(data map[string][]string) *records.Item {
	if data == nil {
		return nil
	}

	var tagItems []tags.Item
	for _, t := range strings.Split(value(data, "tags"), ";") {
		tagItems = append(tagItems, tags.Item(strings.TrimSpace(t)))
	}

	r := records.Item{
		ID:       records.ID(value(data, "id")),
		IssuedID: ns.ID(value(data, "issued_id")),
		Content: records.Content{
			Title:   value(data, "title"),
			Summary: value(data, "summary"),
			// TypeKey:  "",
			Data: value(data, "content_data"),
			// Embedded: nil,
			Tags: tagItems,
		},
		//OwnerID:   "",
		//ViewerID:  "",
	}

	return &r
}

func DataFromRecord(r *records.Item) map[string]string {
	if r == nil {
		return nil
	}

	var updatedAt string
	if r.UpdatedAt != nil {
		updatedAt = r.UpdatedAt.Format("02.01.2006 15:04:05")
	}

	data := map[string]string{
		"id":           string(r.ID),
		"issued_id":    string(r.IssuedID),
		"data_type":    "record", // TODO!!!
		"title":        r.Content.Title,
		"summary":      r.Content.Summary,
		"content_data": r.Content.Data,
		"tags":         strings.Join(r.Content.Tags, "; "),
		// "embedded": r.Content.Embedded,
		"created_at": r.CreatedAt.Format("02.01.2006 15:04:05"),
		"updated_at": updatedAt,
	}

	//linksList, err := json.Marshal(r.Links)
	//if err != nil {
	//	return nil, nil, errors.Wrapf(err, "can't marshal object.tags: %#v for object.id: %s", r.Links, r.ID)
	//}
	//DataFromRecord["links"] = string(linksList)
	//
	//tags := ""
	//filesList := []interfaces.Link{}
	//for _, l := range r.Links {
	//	switch l.Type {
	//
	//	case links.TypeTag:
	//		tags += l.Name + "; "
	//
	//	case files.LinkType:
	//		filesList = append(filesList, l)
	//	}
	//}
	//DataFromRecord["tags"] = tags
	//if len(filesList) > 0 {
	//	files, err := json.Marshal(filesList)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	DataFromRecord["files"] = string(files)
	//}
	//
	//if r.UpdatedAt != nil {
	//	DataFromRecord["updated_at"] = r.UpdatedAt.Format("02.01.2006 15:04:05")
	//}

	return data

}
