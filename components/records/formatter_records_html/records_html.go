package formatter_records_html

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/tools/components/formatter"
	"github.com/pavlo67/tools/components/records"
)

var _ formatter.Operator = &formatterRecordsHTML{}

type formatterRecordsHTML struct{}

const onNew = "on formatterRecordsHTML.New(): "

func New() (formatter.Operator, error) {
	return &formatterRecordsHTML{}, nil
}

const onPrepare = "on formatterRecordsHTML.Prepare(): "

// should be thread-safe
func (dataOp *formatterRecordsHTML) Prepare(key formatter.Key, template string, params common.Map) error {
	return nil
}

const onFormat = "on formatterRecordsHTML.Format(): "

// should be thread-safe
func (dataOp *formatterRecordsHTML) Format(value interface{}, key formatter.Key) (string, error) {
	var recordsItem *records.Item

	switch v := value.(type) {
	case records.Item:
		recordsItem = &v
	case *records.Item:
		if v != nil {
			recordsItem = v
		}
	}

	if recordsItem == nil {
		return "", errors.Errorf("wrong records.Item in %#v", value)
	}

	switch key {
	case formatter.Full:
		return Full(recordsItem), nil
	case formatter.Brief:
		return Brief(recordsItem), nil
	case formatter.Edit:
		return Edit(recordsItem), nil
	}

	return "", fmt.Errorf(onFormat+": wrong formatter.Key (%s)", key)
}
