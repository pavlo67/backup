package records_sqlite

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/connect"
	"github.com/pavlo67/tools/components/records"
)

func Starter() starter.Operator {
	return &recordsSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsSQLiteStarter{}

type recordsSQLiteStarter struct {
	table string

	connectKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (rss *recordsSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rss *recordsSQLiteStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	rss.table, _ = options.String("table")
	rss.connectKey = joiner.InterfaceKey(options.StringDefault("connect_key", string(connect.InterfaceSQLiteKey)))

	rss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records.InterfaceKey)))
	rss.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(records.InterfaceCleanerKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (rss *recordsSQLiteStarter) Run(joinerOp joiner.Operator) error {
	db, _ := joinerOp.Interface(rss.connectKey).(*sql.DB)
	if db == nil {
		return errors.Errorf("no *sql.DB with key %s", rss.connectKey)
	}
	recordsOp, recordsCleanerOp, err := New(db, rss.table, rss.interfaceKey)
	if err != nil {
		return errata.CommonError(err, "can't init records.Operator")
	}

	if err = joinerOp.Join(recordsOp, rss.interfaceKey); err != nil {
		return errata.CommonError(err, fmt.Sprintf("can't join *recordsSQLite as records.Operator with key '%s'", rss.interfaceKey))
	}

	if err = joinerOp.Join(recordsCleanerOp, rss.cleanerKey); err != nil {
		return errata.CommonError(err, fmt.Sprintf("can't join *recordsSQLite as crud.Cleaner with key '%s'", rss.cleanerKey))
	}

	return nil
}
