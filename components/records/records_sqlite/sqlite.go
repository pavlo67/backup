package records_sqlite

// TODO!!! fix according to new data structures

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/selectors/selectors_sql"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/strlib"

	"github.com/pavlo67/tools/components/records"
)

var fieldsToInsert = []string{"title", "summary", "type_key", "data", "embedded", "tags", "issued_id", "owner_id", "viewer_id", "history"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToUpdate = append(fieldsToInsert, "updated_at")
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToRead = append(fieldsToUpdate, "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ records.Operator = &recordsSQLite{}

type recordsSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove, stmClean *sql.Stmt
}

const onNew = "on recordsSQLite.New(): "

func New(db *sql.DB, table string) (records.Operator, crud.Cleaner, error) {
	if table == "" {
		table = records.CollectionDefault
	}

	recordsOp := recordsSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		sqlRemove: "DELETE FROM " + table + " where id = ?",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",

		sqlClean: "DELETE FROM " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&recordsOp.stmInsert, recordsOp.sqlInsert},
		{&recordsOp.stmUpdate, recordsOp.sqlUpdate},
		{&recordsOp.stmRead, recordsOp.sqlRead},
		{&recordsOp.stmRemove, recordsOp.sqlRemove},
		{&recordsOp.stmClean, recordsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &recordsOp, &recordsOp, nil
}

const onSave = "on recordsSQLite.Save(): "

func (dataOp *recordsSQLite) Save(item records.Item, options *crud.Options) (*records.Item, error) {

	if options == nil || options.Identity == nil {
		return nil, errata.CommonError(errata.NoRightsKey)
	}

	// TODO!!! rbac check

	if item.ID == "" {
		// TODO!!!
		item.OwnerID = options.Identity.ID
	}

	var err error

	var embeddedBytes []byte
	if len(item.Content.Embedded) > 0 {
		if embeddedBytes, err = json.Marshal(item.Content.Embedded); err != nil {
			return nil, errors.Wrapf(err, onSave+"can't marshal .Embedded(%#v)", item.Content.Embedded)
		}
	}

	var tagsBytes []byte
	if len(item.Content.Tags) > 0 {
		if tagsBytes, err = json.Marshal(item.Content.Tags); err != nil {
			return nil, errors.Wrapf(err, onSave+"can't marshal .Tags(%#v)", item.Content.Tags)
		}
	}

	// TODO!!! append to .History

	var historyBytes []byte
	if len(item.History) > 0 {
		historyBytes, err = json.Marshal(item.History)
		if err != nil {
			return nil, errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item.History)
		}
	}

	// "title", "summary", "type_key", "data", "embedded", "tags",
	// "issued_id", "owner_id", "viewer_id", "history"
	values := []interface{}{
		item.Content.Title, item.Content.Summary, item.Content.TypeKey, item.Content.Data, embeddedBytes, tagsBytes,
		item.IssuedID, item.OwnerID, item.ViewerID, historyBytes}

	if item.ID == "" {

		res, err := dataOp.stmInsert.Exec(values...)
		if err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, strlib.Stringify(values))
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlInsert, strlib.Stringify(values))
		}
		item.ID = records.ID(strconv.FormatInt(idSQLite, 10))

	} else {
		values = append(values, time.Now().Format(time.RFC3339), item.ID)
		if _, err := dataOp.stmUpdate.Exec(values...); err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlUpdate, strlib.Stringify(values))
		}

	}

	return &item, nil
}

const onRead = "on recordsSQLite.Read(): "

func (dataOp *recordsSQLite) Read(id records.ID, options *crud.Options) (*records.Item, error) {
	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong id (%s)", id)
	}

	item := records.Item{ID: id}

	var embeddedBytes, tagsBytes, historyBytes []byte

	// "title", "summary", "type_key", "data", "embedded", "tags",
	// "issued_id", "owner_id", "viewer_id", "history", "updated_at", "created_at"

	if err = dataOp.stmRead.QueryRow(idNum).Scan(
		&item.Content.Title, &item.Content.Summary, &item.Content.TypeKey, &item.Content.Data, &embeddedBytes, &tagsBytes,
		&item.IssuedID, &item.OwnerID, &item.ViewerID, &historyBytes, &item.UpdatedAt, &item.CreatedAt); err == sql.ErrNoRows {
		return nil, errata.NotFound
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, idNum)
	}

	if len(embeddedBytes) > 0 {
		if err = json.Unmarshal(embeddedBytes, &item.Content.Embedded); err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embeddedBytes)
		}
	}

	if len(tagsBytes) > 0 {
		if err = json.Unmarshal(tagsBytes, &item.Content.Tags); err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tagsBytes)
		}
	}

	if len(historyBytes) > 0 {
		if err = json.Unmarshal(historyBytes, &item.History); err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", historyBytes)
		}
	}

	return &item, nil
}

const onRemove = "on recordsSQLite.Remove()"

func (dataOp *recordsSQLite) Remove(id records.ID, options *crud.Options) error {

	// TODO!!! rbac check

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onRemove+"wrong id (%s)", id)
	}

	if _, err = dataOp.stmRemove.Exec(idNum); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, dataOp.sqlRemove, idNum)
	}

	return nil
}

const onList = "on recordsSQLite.List()"

func (dataOp *recordsSQLite) List(options *crud.Options) ([]records.Item, error) {
	selector := options.GetSelector()
	condition, values, err := selectors_sql.Use(selector)
	if err != nil {
		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", selector, err)
	}

	query := sqllib.SQLList(dataOp.table, fieldsToListStr, condition, options)
	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	}

	//l.Infof("%s / %#v\n%s", condition, values, query)

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []records.Item

	for rows.Next() {
		var idNum int64
		var item records.Item
		var embeddedBytes, tagsBytes, historyBytes []byte

		// "title", "summary", "type_key", "data", "embedded", "tags",
		// "issued_id", "owner_id", "viewer_id", "history", "updated_at", "created_at",
		// "id"

		if err := rows.Scan(
			&item.Content.Title, &item.Content.Summary, &item.Content.TypeKey, &item.Content.Data, &embeddedBytes, &tagsBytes,
			&item.IssuedID, &item.OwnerID, &item.ViewerID, &historyBytes, &item.UpdatedAt, &item.CreatedAt,
			&idNum); err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if len(embeddedBytes) > 0 {
			if err = json.Unmarshal(embeddedBytes, &item.Content.Embedded); err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embeddedBytes)
			}
		}

		if len(tagsBytes) > 0 {
			if err = json.Unmarshal(tagsBytes, &item.Content.Tags); err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tagsBytes)
			}
		}

		if len(historyBytes) > 0 {
			if err = json.Unmarshal(historyBytes, &item.History); err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", historyBytes)
			}
		}

		item.ID = records.ID(strconv.FormatInt(idNum, 10))
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onCount = "on recordsSQLite.Count(): "

func (dataOp *recordsSQLite) Stat(*crud.Options) (common.Map, error) {
	//condition, values, err := selectors_sql.Use(term)
	//if err != nil {
	//	termStr, _ := json.Marshal(term)
	//	return 0, errors.Wrapf(err, onCount+": can't selectors_sql.Use(%s)", termStr)
	//}
	//
	//query := sqllib.SQLCount(dataOp.table, condition, options)
	//stm, err := dataOp.db.Prepare(query)
	//if err != nil {
	//	return 0, errors.Wrapf(err, onCount+": can't db.Prepare(%s)", query)
	//}
	//
	//var num uint64
	//
	//err = stm.QueryRow(values...).Scan(&num)
	//if err != nil {
	//	return 0, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, query, values)
	//}

	return nil, errata.NotImplemented
}

func (dataOp *recordsSQLite) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on recordsSQLite.Close()")
}

//
//const onDetails = "on recordsSQLite.Details(): "
//
//func (dataOp *recordsSQLite) SetDetails(item *records.Item) error {
//	if item == nil {
//		return errors.New(onDetails + "nil item")
//	}
//
//	// l.Infof("11111111111 %s %s %t", item.DetailsRaw, item.TypeKey, item.TypeKey == records.TypeKeyTest)
//
//	if len(item.DetailsRaw) < 1 {
//		item.Details = nil
//		return nil
//	}
//
//	switch item.TypeKey {
//	case types.KeyString:
//		item.Details = string(item.DetailsRaw)
//
//	case records.TypeKeyTest:
//		item.Details = &records.Test{}
//		err := json.Unmarshal(item.DetailsRaw, item.Details)
//		if err != nil {
//			return errors.Wrapf(err, onDetails+"can't .Unmarshal(%#v)", item.DetailsRaw)
//		}
//
//	default:
//
//		// TODO: remove the kostyl
//		item.Details = string(item.DetailsRaw)
//
//		// return errors.Errorf(onDetails+"unknown item.TypeKey(%s) for item.DetailsRaw(%s)", item.TypeKey, item.DetailsRaw)
//
//	}
//
//	// l.Infof("11111111111 %#v", item.Details)
//
//	return nil
//}
//
//const onExport = "on recordsSQLite.Export()"
//
//func (dataOp *recordsSQLite) Export(afterIDStr string, options *crud.Options) ([]records.Item, error) {
//	// TODO: remove limits
//	// if options != nil {
//	//	options.Limits = nil
//	// }
//
//	afterIDStr = strings.TrimSpace(afterIDStr)
//
//	var term *selectors.Term
//
//	var afterID int
//	if afterIDStr != "" {
//		var err error
//		afterID, err = strconv.Atoi(afterIDStr)
//		if err != nil {
//			return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter: %s", afterIDStr, err)
//		}
//
//		// TODO!!! term with some item's autoincrement if original .Key isn't it (using .Key to find corresponding autoincrement value)
//		term = selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
//	}
//
//	// TODO!!! order by some item's autoincrement if original .Key isn't it
//	if options == nil {
//		options = &crud.Options{OrderBy: []string{"id"}}
//	} else {
//		options.OrderBy = []string{"id"}
//	}
//
//	return dataOp.List(term, options)
//}
