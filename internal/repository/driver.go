package repository

import (
	"fmt"

	"github.com/slovaq/web2tg/internal/repository/sqlite3"
)

type iDb interface {
	Connect(url string, dbName string) error
	Disconnect() error
	DriverType() string
	GetDriverImplementation() interface{}
	Migrate(tableName string, Schema interface{}) error
}

func CreateDriver(dType string) (iDb, error) {
	switch dType {
	// case "mongo":
	// 	return mongolib.NewDriver()
	case "sqlite3":
		return sqlite3.NewDriver(), nil
	}
	return nil, fmt.Errorf("driver: driver type not supported or config is broken")
}
