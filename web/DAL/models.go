package DAL

import (
	"time"
)

type City struct {
	ID     uint   `gorm:"primarykey"`
	Name   string `gorm:"unique"`
	ChatID int
}
type Record struct {
	ID          uint `gorm:"primarykey"`
	Date        time.Time
	Destination []City `gorm:"many2many:record_destinations;"`
	Text        string
}
type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Login    string `gorm:"unique"`
	Password string
}
type ClientConfig struct {
	Login    string
	ChatLink string
	BotToken string
}
