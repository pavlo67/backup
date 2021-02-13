package records_sqlite

import (
	"errors"
	"strings"

	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/tools/components/records"
)

const onHasTag = "on recordsSQLite.HasTag()"

func (recordsOp *recordsSQLite) HasTag(tag string) (selectors.Term, error) {
	if tag = strings.TrimSpace(tag); tag == "" {
		return nil, errors.New(onHasTag + ": no tag to select records")
	}

	return selectors.TermSQL{Condition: `tags LIKE ?`, Values: []interface{}{`%"` + tag + `"%`}}, nil
}

const onAddParent = "on recordsSQLite.AddParent()"

func (recordsOp *recordsSQLite) AddParent(tags []string, id records.ID) ([]string, error) {
	if id = records.ID(strings.TrimSpace(string(id))); id == "" {
		return nil, errors.New(onAddParent + ": no id to add parent record")
	}

	return append(tags, string(id)+":"), nil
}

const onHasParent = "on recordsSQLite.HasParent()"

func (recordsOp *recordsSQLite) HasParent(id records.ID) (selectors.Term, error) {
	if id = records.ID(strings.TrimSpace(string(id))); id == "" {
		return nil, errors.New(onHasParent + ": no id to select child records")
	}

	return selectors.TermSQL{Condition: `tags LIKE ?`, Values: []interface{}{`%"` + string(id) + `:%`}}, nil
}
