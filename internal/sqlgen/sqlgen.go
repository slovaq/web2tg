package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
)

type Column struct {
	Name  string
	Type  string
	Index bool
}

func CreateTableGen(name string, table interface{}) (string, error) {
	typ := reflect.TypeOf(table)
	//fmt.Printf("typ: %v\n", typ)
	//	fmt.Printf("kind: %v\n", typ.Kind())
	if typ.Kind() != reflect.Struct {
		return "", fmt.Errorf("v must be a struct")
	}
	var columns []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		columnName := field.Tag.Get("field")
		if columnName == "" {
			columnName = strings.ToLower(field.Name)
		}
		sqlType, err := getSQLType(field.Type)
		if err != nil {
			return "", err
		}
		columns = append(columns, columnName+" "+sqlType)
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", name, strings.Join(columns, ", "))
	return query, nil
}

func contains(s []Column, e Column) bool {
	for _, a := range s {
		if a.Name == e.Name {
			return true
		}
	}
	return false
}
func AlterTableGen(name string, table interface{}, alreadyExist []Column) (string, bool, error) {
	typ := reflect.TypeOf(table)
	if typ.Kind() != reflect.Struct {
		return "", false, fmt.Errorf("v must be a struct")
	}
	//var columns []string
	var newColumns []Column
	colsql := ""
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		columnName := field.Tag.Get("field")
		if columnName == "" {
			columnName = strings.ToLower(field.Name)
		}
		sqlType, err := getSQLType(field.Type)
		if err != nil {
			return "", false, err
		}
		newColumn := Column{
			Name: columnName,
			Type: sqlType,
		}
		if !contains(alreadyExist, newColumn) {
			newColumns = append(newColumns, newColumn)
			colsql = colsql + fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", name, columnName, sqlType)
		}
		//columns = append(columns, columnName+" "+sqlType)
		//
	}
	if len(newColumns) == 0 {
		return "", false, nil
	}
	//	query := fmt.Sprintf("ALTER TABLE %s\n %s", name, colsql)
	return colsql, true, nil
}

// func AddIndexGen
func getSQLType(typ reflect.Type) (string, error) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER", nil
	case reflect.Float32, reflect.Float64:
		return "REAL", nil
	case reflect.Bool:
		return "INTEGER", nil
	case reflect.String:
		return "TEXT", nil
	default:
		return "", fmt.Errorf("unsupported type: %v", typ)
	}
}
