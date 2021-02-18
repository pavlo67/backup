package records_sqlite

import (
	"errors"
	"strings"

	"github.com/pavlo67/tools/components/tags"

	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/tools/components/records"
)

const onHasTag = "on recordsSQLite.HasTag()"

func (recordsOp *recordsSQLite) HasTag(tag tags.Item) (selectors.Term, error) {
	tagStr := strings.TrimSpace(string(tag))
	if tagStr == "" {
		return nil, errors.New(onHasTag + ": no tagStr to select records")
	}

	return selectors.TermSQL{Condition: `tags LIKE ?`, Values: []interface{}{`%"` + tagStr + `"%`}}, nil
}

const onHasNoTag = "on recordsSQLite.HasNoTag()"

func (recordsOp *recordsSQLite) HasNoTag() (selectors.Term, error) {
	return selectors.TermSQL{Condition: `tags IN ('', '{}')`, Values: nil}, nil
}

const onAddParent = "on recordsSQLite.AddParent()"

func (recordsOp *recordsSQLite) AddParent(ts []tags.Item, id records.ID) ([]tags.Item, error) {
	idStr := strings.TrimSpace(string(id))
	if idStr == "" {
		return nil, errors.New(onAddParent + ": no id to add parent record")
	}

	return append(ts, tags.Item(idStr+":")), nil
}

const onHasParent = "on recordsSQLite.HasParent()"

func (recordsOp *recordsSQLite) HasParent(id records.ID) (selectors.Term, error) {
	if id = records.ID(strings.TrimSpace(string(id))); id == "" {
		return nil, errors.New(onHasParent + ": no id to select child records")
	}

	return selectors.TermSQL{Condition: `tags LIKE ?`, Values: []interface{}{`%"` + string(id) + `:%`}}, nil
}
