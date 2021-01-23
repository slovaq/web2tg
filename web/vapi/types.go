package vapi

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

type Link struct {
	ID       int `gorm:"primaryKey"`
	UserLink string
	ChatID   int64 `gorm:"unique"`
}

type ClientConfig struct {
	Login    string
	City     string
	ChatLink string
	BotToken string
}

type VapiRecord struct {
	User    string
	Message string
	City    string
	Date    string
	ID      int `gorm:"primarykey"`
	Time    string
	Status  string
	Period  string
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

//Config struct {Token string, UpdateTime int}
type Config struct {
	Token      string
	UpdateTime int
}

//SNBot struct { cfg *Config, bot *tgbotapi.BotAPI, upd tgbotapi.UpdatesChannel}
type SNBot struct {
	cfg *Config
	bot *tgbotapi.BotAPI
	upd tgbotapi.UpdatesChannel
}
