package DAL

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
