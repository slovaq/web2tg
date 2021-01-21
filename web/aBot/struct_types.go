package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

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
