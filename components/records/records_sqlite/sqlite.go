package records_sqlite

// TODO!!! fix according to new data structures

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/strlib"

	"github.com/pavlo67/tools/components/records"
	"github.com/pavlo67/tools/components/tags"
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

func (recordsOp *recordsSQLite) Save(item records.Item, options *crud.Options) (*records.Item, error) {

	//if options == nil || options.Identity == nil {
	//	return nil, errors.CommonError(common.NoRightsKey)
	//}

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

		res, err := recordsOp.stmInsert.Exec(values...)
		if err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantExec, recordsOp.sqlInsert, strlib.Stringify(values))
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, recordsOp.sqlInsert, strlib.Stringify(values))
		}
		item.ID = records.ID(strconv.FormatInt(idSQLite, 10))

	} else {
		values = append(values, time.Now().Format(time.RFC3339), item.ID)
		if _, err := recordsOp.stmUpdate.Exec(values...); err != nil {
			return nil, errors.Wrapf(err, onSave+sqllib.CantExec, recordsOp.sqlUpdate, strlib.Stringify(values))
		}

	}

	return &item, nil
}

const onRead = "on recordsSQLite.Read(): "

func (recordsOp *recordsSQLite) Read(id records.ID, options *crud.Options) (*records.Item, error) {
	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(onRead+"wrong id (%s)", id)
	}

	item := records.Item{ID: id}

	var embeddedBytes, tagsBytes, historyBytes []byte

	// "title", "summary", "type_key", "data", "embedded", "tags",
	// "issued_id", "owner_id", "viewer_id", "history", "updated_at", "created_at"

	if err = recordsOp.stmRead.QueryRow(idNum).Scan(
		&item.Content.Title, &item.Content.Summary, &item.Content.TypeKey, &item.Content.Data, &embeddedBytes, &tagsBytes,
		&item.IssuedID, &item.OwnerID, &item.ViewerID, &historyBytes, &item.UpdatedAt, &item.CreatedAt); err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, recordsOp.sqlRead, idNum)
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

func (recordsOp *recordsSQLite) Remove(id records.ID, options *crud.Options) error {

	// TODO!!! rbac check

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return fmt.Errorf(onRemove+"wrong id (%s)", id)
	}

	if _, err = recordsOp.stmRemove.Exec(idNum); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, recordsOp.sqlRemove, idNum)
	}

	return nil
}

const onList = "on recordsSQLite.List()"

func (recordsOp *recordsSQLite) List(options *crud.Options) ([]records.Item, error) {

	var termSQL selectors.TermSQL

	if selector := options.GetSelector(); selector != nil {
		var ok bool
		if termSQL, ok = selector.(selectors.TermSQL); !ok {
			return nil, fmt.Errorf(onList+": wrong selector: %#v", selector)
		}
	}

	query := sqllib.SQLList(recordsOp.table, fieldsToListStr, termSQL.Condition, options)
	stm, err := recordsOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	}

	//l.Infof("%s / %#v\n%s", condition, values, query)

	rows, err := stm.Query(termSQL.Values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, query, termSQL.Values)
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
			return items, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, query, termSQL.Values)
		}

		if len(embeddedBytes) > 0 {
			if err = json.Unmarshal(embeddedBytes, &item.Content.Embedded); err != nil {
				return items, errors.Wrapf(err, onList+": can't unmarshal .Embedded (%s)", embeddedBytes)
			}
		}

		if len(tagsBytes) > 0 {
			if err = json.Unmarshal(tagsBytes, &item.Content.Tags); err != nil {
				return items, errors.Wrapf(err, onList+": can't unmarshal .Tags (%s)", tagsBytes)
			}
		}

		if len(historyBytes) > 0 {
			if err = json.Unmarshal(historyBytes, &item.History); err != nil {
				return items, errors.Wrapf(err, onList+": can't unmarshal .History (%s)", historyBytes)
			}
		}

		item.ID = records.ID(strconv.FormatInt(idNum, 10))
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, termSQL.Values)
	}

	return items, nil
}

const onTags = "on recordsSQLite.Tags()"

func (recordsOp *recordsSQLite) Tags(options *crud.Options) (tags.StatMap, error) {

	var termSQL selectors.TermSQL

	if selector := options.GetSelector(); selector != nil {
		var ok bool
		if termSQL, ok = selector.(selectors.TermSQL); !ok {
			return nil, fmt.Errorf(onTags+": wrong selector: %#v", selector)
		}
	}

	query := sqllib.SQLList(recordsOp.table, "tags", termSQL.Condition, options)
	stm, err := recordsOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onTags+": can't db.Prepare(%s)", query)
	}

	rows, err := stm.Query(termSQL.Values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onTags+": "+sqllib.CantQuery, query, termSQL.Values)
	}
	defer rows.Close()

	tagsStat := tags.StatMap{}

	for rows.Next() {
		var tagsBytes []byte

		if err := rows.Scan(&tagsBytes); err != nil {
			return tagsStat, errors.Wrapf(err, onTags+": "+sqllib.CantScanQueryRow, query, termSQL.Values)
		}

		if len(tagsBytes) > 0 {
			var ts []tags.Item
			if err = json.Unmarshal(tagsBytes, &ts); err != nil {
				// TODO!!! collect errors
				l.Errorf(onTags+": can't unmarshal ts (%s): %s", tagsBytes, err)
				continue
			}

			for _, tag := range ts {
				// l.Info(tag, tagsStat[tag])

				// ts := tagsStat[tag]
				tagsStat[tag] = tagsStat[tag] + 1
			}
		}

	}

	if err = rows.Err(); err != nil {
		return tagsStat, errors.Wrapf(err, onTags+": "+sqllib.RowsError, query, termSQL.Values)
	}

	return tagsStat, nil
}

func (recordsOp *recordsSQLite) Close() error {
	return errors.Wrap(recordsOp.db.Close(), "on recordsSQLite.Close()")
}

//const onStat = "on recordsSQLite.StatMap(): "
//
//func (recordsOp *recordsSQLite) StatMap(*crud.Options) (common.Map, error) {
//	condition, values, err := selectors_sql.Use(term)
//	if err != nil {
//		termStr, _ := json.Marshal(term)
//		return 0, errors.Wrapf(err, onCount+": can't selectors_sql.Use(%s)", termStr)
//	}
//
//	query := sqllib.SQLCount(recordsOp.table, condition, options)
//	stm, err := recordsOp.db.Prepare(query)
//	if err != nil {
//		return 0, errors.Wrapf(err, onCount+": can't db.Prepare(%s)", query)
//	}
//
//	var num uint64
//
//	err = stm.QueryRow(values...).Scan(&num)
//	if err != nil {
//		return 0, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, query, values)
//	}
//
//	return nil, common.ErrNotImplemented
//}
