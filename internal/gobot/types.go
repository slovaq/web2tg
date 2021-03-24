package gobot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MessageTG struct {
	Message string
	ChatID  string
	Pic     string
}
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
type GobotConnect struct {
	MessageTG   chan MessageTG
	Updatetoken chan bool
	CheckInit   chan bool
}

//SNBot struct { cfg *Config, bot *tgbotapi.BotAPI, upd tgbotapi.UpdatesChannel}
type SNBot struct {
	cfg *Config
	bot *tgbotapi.BotAPI
	upd tgbotapi.UpdatesChannel
}
type UpdateStorage struct {
	UpdateRecord chan bool
	UpdateConfig chan string
	ReadRecord   chan bool
	ReadConfig   chan string

	Box []Box
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

//Config struct {Token string, UpdateTime int}
type Config struct {
	Token      string
	UpdateTime int
}
