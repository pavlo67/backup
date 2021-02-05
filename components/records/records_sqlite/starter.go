package records_sqlite

import (
	"github.com/pavlo67/tools/components/records"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &recordsSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsSQLiteStarter{}

type recordsSQLiteStarter struct {
	config config.Access
	table  string

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (ts *recordsSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *recordsSQLiteStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgSQLite config.Access
	err := cfg.Value("sqlite", &cfgSQLite)
	if err != nil {
		return nil, err
	}

	ts.config = cfgSQLite
	ts.table, _ = options.String("table")
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records.InterfaceKey)))
	ts.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(records.CleanerInterfaceKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (ts *recordsSQLiteStarter) Setup() error {
	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Title: sm.table}},
	//)
}

func (ts *recordsSQLiteStarter) Run(joinerOp joiner.Operator) error {
	var ok bool

	//if !ts.noTagger {
	//	taggerOp, ok = joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
	//	if !ok {
	//		return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
	//	}
	//
	//	taggercleanerOp, ok = joinerOp.Interface(tagger.CleanerInterfaceKey).(crud.Cleaner)
	//	if !ok {
	//		return errors.Errorf("no tagger.Cleaner with key %s", tagger.InterfaceKey)
	//	}
	//}
	//
	recordsOp, recordscleanerOp, err := New(ts.config, ts.table, ts.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init records.Operator")
	}

	err = joinerOp.Join(recordsOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *recordsSQLite as records.Operator with key '%s'", ts.interfaceKey)
	}

	err = joinerOp.Join(recordscleanerOp, ts.cleanerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *recordsSQLite as crud.Cleaner with key '%s'", ts.cleanerKey)
	}

	return nil
}
