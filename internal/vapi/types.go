package vapi

import (
	"github.com/slovaq/web2tg/internal/gobot"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	//Logger: logger.Default.LogMode(logger.Silent)
})

type Data struct {
	User string
}

type Link struct {
	ID       int `gorm:"primaryKey"`
	UserLink string
	ChatID   int64 `gorm:"unique"`
}

//PostSorter []VapiRecord
type PostSorter []VapiRecord
type ClientConfig struct {
	Login    string
	City     string
	ChatLink string
	BotToken string
}

type VapiRecord struct {
	User     string
	Message  string
	City     string
	Date     string
	ID       int `gorm:"primarykey"`
	Time     string
	Status   string
	Period   string
	DataRead string
	Pic      string
}
type CreateConfData struct {
	User   *ClientConfig
	Status string
}
type tomlConfig struct {
	Token string
}

//Update struct{UpdateToken chan string}
type Update struct {
	UpdateToken chan string
}

type infoMutex struct {
	Status    bool
	Locker    string
	MutexName string
}

//Box Message string, Time int64, Token string, URL string, ID int, User string
type Box struct {
	Message string
	Time    int64
	Token   string
	URL     string
	ID      int
	User    string
	Status  string
	Period  string
	Pic     string
}

//Boxs []Box
type Boxs []Box

//IntRange struct {min, max int}
type IntRange struct {
	min, max int
}
type UpdateStorage struct {
	UpdateRecord chan bool
	UpdateConfig chan string
	ReadRecord   chan bool
	ReadConfig   chan string
	Box          []Box
	GobotConnect gobot.GobotConnect
}
type HandleError struct {
	HttpError string
}
