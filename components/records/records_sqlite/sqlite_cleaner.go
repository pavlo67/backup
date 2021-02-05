package records_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/sqllib"
)

var _ crud.Cleaner = &dataSQLite{}

const onIDs = "on dataSQLite.IDs()"

//func (dataOp *dataSQLite) ids(condition string, values []interface{}) ([]interface{}, error) {
//	if strings.TrimSpace(condition) != "" {
//		condition = " WHERE " + condition
//	}
//
//	query := "SELECT id FROM " + dataOp.table + condition
//	stm, err := dataOp.db.Prepare(query)
//	if err != nil {
//		return nil, errors.Wrapf(err, onIDs+": can't db.Prepare(%s)", query)
//	}
//
//	rows, err := stm.Query(values...)
//	if err == sql.ErrNoRows {
//		return nil, nil
//	} else if err != nil {
//		return nil, errors.Wrapf(err, onIDs+sqllib.CantQuery, query, values)
//	}
//	defer rows.Close()
//
//	var ids []interface{}
//
//	for rows.Next() {
//		var id common.ID
//
//		err := rows.Scan(&id)
//		if err != nil {
//			return ids, errors.Wrapf(err, onIDs+sqllib.CantScanQueryRow, query, values)
//		}
//
//		ids = append(ids, id)
//	}
//	err = rows.Err()
//	if err != nil {
//		return ids, errors.Wrapf(err, onIDs+": "+sqllib.RowsError, query, values)
//	}
//
//	return ids, nil
//}

const onClean = "on dataSQLite.Clean(): "

func (dataOp *dataSQLite) Clean(_ *crud.Options) error {
	//var termTags *selectors.Term
	//
	//condition, values, err := selectors_sql.Use(nil)
	//if err != nil {
	//	return errors.Errorf(onClean+"wrong selector: %s", err)
	//}

	query := dataOp.sqlClean

	//if strings.TrimSpace(condition) != "" {
	//	ids, err := dataOp.ids(condition, values)
	//	if err != nil {
	//		return errors.Wrap(err, onClean+"can't dataOp.ids(condition, values)")
	//	}
	//	termTags = logic.AND(selectors.In("key", dataOp.interfaceKey), selectors.In("id", ids...))
	//
	//	query += " WHERE " + condition
	//
	//} else {
	//	termTags = selectors.In("key", dataOp.interfaceKey) // TODO!!! correct field key
	//
	//}

	var values []interface{}
	_, err := dataOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	//if dataOp.taggerCleaner != nil {
	//	err = dataOp.taggerCleaner.Clean(termTags, nil)
	//	if err != nil {
	//		return errors.Wrap(err, onClean)
	//	}
	//}

	return err
}

//const onSelectToClean = "on dataSQLite.SelectToClean(): "
//
//func (dataOp *dataSQLite) SelectToClean(options *crud.Options) (*selectors.Term, error) {
//	var limit uint64
//
//	if options != nil && options.Limit > 0 {
//		limit = options.Limit
//	} else {
//		return nil, errors.New(onSelectToClean + "no clean limit is defined")
//	}
//
//	queryMax := "SELECT MAX(id) from " + dataOp.table
//
//	var maxID uint64
//	row := dataOp.db.QueryRow(queryMax)
//
//	err := row.Scan(&maxID)
//	if err != nil {
//		return nil, errors.Errorf(onSelectToClean+": error on query (%s)", queryMax)
//	}
//
//	return selectors.Binary(selectors.Le, "id", selectors.Value{V: maxID - limit}), nil
//}

//queryDelete := "DELETE from " + dataOp.table + " WHERE id <= ?"
//res, err := dataOp.db.Exec(queryDelete, maxID-limit)
//if err != nil {
//return errors.Errorf(onClean+": error on query (%s)", queryDelete)
//}
//
//rowsAffected, err := res.RowsAffected()
//if err != nil {
//return errors.Errorf(onClean+": error on res.RowsAffected(%s)", queryDelete)
//}
//
//l.Infof(onClean+": res.RowsAffected() = %d", rowsAffected)
//
//if dataOp.tableTags != "" {
//// TODO!!!
//}
//
