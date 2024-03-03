package efactory

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// SQLDB is for raw SQL database operations
type SQLDB struct {
	// DB is the database connection
	// must provide if want to insert data into the database
	DB *sql.DB

	// Ctx is the context for the database operations
	// it is optional
	Ctx context.Context
}

type inserParams struct {
	tableName string
	value     interface{}
}

type inserListParams struct {
	tableName string
	values    []interface{}
}

// Implement the insert method for SQLDB
func (s *SQLDB) insert(params inserParams) (interface{}, error) {
	rawStmt, vals := prepareStmtAndVals(params.tableName, params.value)

	// Prepare the insert statement
	stmt, err := s.DB.Prepare(rawStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id, err := insertToDB(s.Ctx, tx, stmt, vals[0])
	if err != nil {
		return nil, err
	}

	setIDField(params.value, id)
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return params.value, nil
}

func (s *SQLDB) insertList(params inserListParams) ([]interface{}, error) {
	rawStmt, fieldValues := prepareStmtAndVals(params.tableName, params.values...)

	stmt, err := s.DB.Prepare(rawStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result := make([]interface{}, len(fieldValues))
	for i, vals := range fieldValues {
		id, err := insertToDB(s.Ctx, tx, stmt, vals)
		if err != nil {
			return nil, err
		}

		v := params.values[i]
		setIDField(v, id)

		result[i] = v
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// prepareStmtAndVals prepares the SQL insert statement and the values to be inserted
// values are the pointer to the struct
func prepareStmtAndVals(tableName string, values ...interface{}) (string, [][]interface{}) {
	fieldNames := []string{}
	placeholders := []string{}
	fieldValues := [][]interface{}{}

	for index, val := range values {
		val := reflect.ValueOf(val).Elem()
		vals := []interface{}{}

		for i := 0; i < val.NumField(); i++ {
			n := val.Type().Field(i).Name
			if n == "ID" {
				continue
			}

			vals = append(vals, val.Field(i).Interface())

			if index == 0 {
				fieldName := val.Type().Field(i).Tag.Get("eDB")
				if fieldName == "" {
					fieldName = camelToSnake(n)
				}

				fieldNames = append(fieldNames, fieldName)
				placeholders = append(placeholders, "?")
			}
		}

		fieldValues = append(fieldValues, vals)
	}

	// Construct the SQL insert statement
	fns := strings.Join(fieldNames, ", ")
	phs := strings.Join(placeholders, ", ")
	rawStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, fns, phs)

	return rawStmt, fieldValues
}

// insertToDB inserts the given values to the database
func insertToDB(ctx context.Context, tx *sql.Tx, stmt *sql.Stmt, vals []interface{}) (int64, error) {
	var res sql.Result
	var errSQL error

	if ctx == nil {
		res, errSQL = tx.Stmt(stmt).Exec(vals...)
	} else {
		res, errSQL = tx.Stmt(stmt).ExecContext(ctx, vals...)
	}

	if errSQL != nil {
		return 0, errSQL
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// setIDField sets the id value on ID field of the given value
func setIDField(v interface{}, id int64) {
	val := reflect.ValueOf(v).Elem()
	idField := val.FieldByName("ID")
	switch idField.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		idField.SetInt(id)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		idField.SetUint(uint64(id))
	}
}
