package notorm

// TODO
// Insert() should take pointer?
// avoid mysql keyword

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func splitTag(tag string) map[string]string {
	ret := make(map[string]string)
	list := strings.Split(tag, ",")
	for _, seg := range list {
		kv := strings.Split(seg, "=")
		if len(kv) != 2 {
			return nil
		}
		ret[kv[0]] = kv[1]
	}

	return ret
}
func sqlType(f reflect.StructField) string {
	tags := splitTag(f.Tag.Get("mysql"))
	typ, ok := tags["type"]
	if ok {
		return typ // TODO: check type
	}

	t := f.Type.Kind()
	switch {
	case t == reflect.String:
		return "VARCHAR(256)"
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
		fmt.Printf("%v", err)
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
		values = append(values, name+" "+sqlType(f))
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table + " (" + strings.Join(values, ", ") + ");"
	if no.debug {
		fmt.Println(sql)
	}
	_, err := no.db.Exec(sql)

	return err
}

func fieldValue(value interface{}, kind reflect.Kind) string {
	switch {
	case kind == reflect.String:
		return `'` + strings.Replace(value.(string), "'", `\'`, -1) + `'`
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
	r, err := no.db.Exec(sql)
	if err != nil {
		fmt.Printf("failed: %v", err)
	}
	rows, err := r.RowsAffected()
	if rows != 1 {
		fmt.Printf("failed: %v", r)
	}
}

// Select a single row and write to a point to a structure
func (no *NotOrm) Select(where string, o interface{}) error {
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

// Select **ALL** rows and write to an array of to a structure
// NOTE: no pagination support
// o should be an empty struct that indicates which table to search
// returns a list of all elements that fit the where string
func (no *NotOrm) SelectAll(where string, o interface{}) ([]interface{}, error) {
	_type := reflect.TypeOf(o)
	sql := "SELECT * FROM " + tableName(_type.Name()) + " " + where + ";"
	if no.debug {
		fmt.Println(sql)
	}
	var arr []interface{}
	rows, err := no.db.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var fields []interface{}
		val := reflect.New(_type) // creates new object each time
		for i := 0; i < _type.NumField(); i++ {
			// There should be a better way of getting the address of a field
			// from a struct pointer.
			fields = append(fields, reflect.Indirect(val).Field(i).Addr().Interface())
		}
		err := rows.Scan(fields...)
		if err != nil {
			break
		}
		arr = append(arr, val.Interface())
	}
	if err != nil {
		return nil, err
	} else {
		return arr, nil
	}
}
