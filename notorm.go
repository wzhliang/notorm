package notorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func sqlType(t reflect.Kind) string {
	switch {
	case t == reflect.String:
		return "VARCHAR"
	case t == reflect.Int:
		return "INTEGER"
	}
	panic("Unknown type")
}

func tableName(name string) string {
	return strings.ToLower(name) + "s"
}

func fieldName(name string) string {
	return strings.ToLower(name)
}

func getField(v interface{}, field string, _type reflect.Kind) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	switch {
	case _type == reflect.String:
		return f.String()
	case _type == reflect.Int:
		return f.Int()
	}
	panic("I don't know how to handle this")
}

type NotOrm struct {
	db    *sql.DB
	debug bool
}

func NewConnection(driver string, param string) *NotOrm {
	no := new(NotOrm)
	db, err := sql.Open(driver, param)
	if err != nil {
		return nil
	} else {
		no.db = db
		return no
	}
}

func (no *NotOrm) Debug(d bool) {
	no.debug = d
}

func (no *NotOrm) CreateTable(o interface{}) error {
	t := reflect.TypeOf(o)
	table := tableName(t.Name())
	var values []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := fieldName(f.Name)
		values = append(values, name+" TYPE "+sqlType(f.Type.Kind()))
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table + " (" + strings.Join(values, ",") + ");"
	if no.debug {
		fmt.Println(sql)
	}
	no.db.Exec(sql)

	return nil
}

func fieldValue(value interface{}, kind reflect.Kind) string {
	switch {
	case kind == reflect.String:
		return `"` + value.(string) + `"`
	case kind == reflect.Int:
		return strconv.FormatInt(value.(int64), 10)
	}
	panic("Unknown type")
}

func (no *NotOrm) Insert(o interface{}) {
	t := reflect.TypeOf(o)
	table := tableName(t.Name())
	var values []string
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, fieldName(f.Name))
		v := getField(o, f.Name, f.Type.Kind())
		values = append(values, fieldValue(v, f.Type.Kind()))
	}
	sql := "INSERT INTO " + table + " (" + strings.Join(fields, ",") + ")" +
		" VALUES (" + strings.Join(values, ",") + ");"
	if no.debug {
		fmt.Println(sql)
	}
	no.db.Exec(sql)
}

func (no *NotOrm) SelectAll(where string, o interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(o)) // o should be a pointer
	_type := val.Type()
	var fields []interface{}
	for i := 0; i < _type.NumField(); i++ {
		fields = append(fields, val.Field(i).Addr().Interface()) // don't miss Interface()
	}
	sql := "SELECT * FROM " + tableName(_type.Name()) + " " + where + " LIMIT 1;"
	if no.debug {
		fmt.Println(sql)
	}
	err := no.db.QueryRow(sql).Scan(fields...)
	if err == nil {
		return err
	}
	return nil
}
