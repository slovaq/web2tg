package vapi

import (
	"sync"

	"github.com/slovaq/web2tg/internal/gobot"
)

var (
	m sync.Mutex
)

func InitChannel(UpdateRecord chan bool, UpdateConfig chan string, ReadConfig chan string, Box Boxs, GobotConnect gobot.GobotConnect) *UpdateStorage {
	return &UpdateStorage{
		UpdateRecord: UpdateRecord,
		UpdateConfig: UpdateConfig,

		ReadConfig:   ReadConfig,
		Box:          Box,
		GobotConnect: GobotConnect,
	}
}
