package DAL

import (
	"gorm.io/gorm"
	"time"
)

type City struct {
	gorm.Model
	Name   string `gorm:"unique"`
	ChatID int
}
type Record struct {
	gorm.Model
	Date        time.Time
	Destination []City `gorm:"many2many:record_destinations;"`
	Text        string
}
type User struct {
	gorm.Model
	Name     string
	Login    string `gorm:"unique"`
	Password string
}
