package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/slovaq/web2tg/internal/sqlgen"
)

type ISqlite3 struct {
	DB *sql.DB
}

func NewDriver() *ISqlite3 {
	return &ISqlite3{}
}

// Connect(url string, dbName string) error
func (i *ISqlite3) Connect(typeDb string, dbName string) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	//ping db
	err = db.Ping()
	if err != nil {
		return err
	}
	i.DB = db
	return nil
}

// Migrate(tableName string,Schema interface{}) error
func (i *ISqlite3) Migrate(tableName string, Schema interface{}) error {

	createSql, err := sqlgen.CreateTableGen(tableName, Schema)
	if err != nil {
		return err
	}
	fmt.Printf("createSql: %v\n", createSql)
	_, err = i.DB.Exec(createSql)
	if err != nil {
		return err
	}
	rows, err := i.DB.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var columns []sqlgen.Column
	for rows.Next() {
		var cid, name, type_ string
		var notnull bool
		var dflt_value, pk interface{}
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return err
		}
		//fmt.Printf("column %s %s %v %v %v\n", name, type_, notnull, dflt_value, pk)
		columns = append(columns, sqlgen.Column{
			Name: name,
			Type: type_,
		})
	}
	// Вывод списка полей
	// for _, column := range columns {
	// 	fmt.Printf("сolumn: %s\n", col)
	// }

	//fmt.Printf("res: %v\n", res)
	alterSql, modified, err := sqlgen.AlterTableGen(tableName, Schema, columns)
	if err != nil {
		return err
	}
	if modified {
		fmt.Printf("alterSql: %v\n", alterSql)
		_, err := i.DB.Exec(alterSql)
		if err != nil {
			return err
		}
		//fmt.Printf("res: %v\n", res)
	}

	return nil
}
func (i *ISqlite3) Disconnect() error {
	if i.DB != nil {
		return i.DB.Close()
	}
	return nil
}
func (i *ISqlite3) DriverType() string {
	return "sqlite3"
}
func (i *ISqlite3) GetDriverImplementation() interface{} {
	return i
}
