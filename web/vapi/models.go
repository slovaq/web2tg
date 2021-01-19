package vapi

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
