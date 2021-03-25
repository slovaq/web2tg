package vapi

import (
	"math/rand"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/data"
	"github.com/slovaq/web2tg/internal/gobot"
)

//Initrc  start sheduler module
func (upd *UpdateStorage) Initrc() {
	stateM.MutexName = "stateM"
	stateM.Status = false
	rand.Seed(time.Now().UnixNano())
	fmc.Printfln("#rbt Run> #gbtInitrc")
	go upd.Check()
	go upd.checkDateCounter()
	go upd.read()
	upd.ReadRecord <- true
	go checkerChannel()
	//go gobot.InitBotRC(&upd.GobotConnect)
	upd.GobotConnect.InitBot()

}

//box.add(2)
//	f := make(chan bool)
//go upd.Check()
//go upd.checkDateCounter()
//sgo upd.read()
func (upd *UpdateStorage) ManageMessage(f Box) {
	fmc.Printfln("#ybtManageMessage> #bbtrun")

	msg := gobot.MessageTG{
		Message: f.Message,
		ChatID:  f.URL,
		Pic:     f.Pic,
	}
	upd.GobotConnect.MessageTG <- msg
	fmc.Printfln("#ybtManageMessage> #bbtsend")
	if f.Period == "one" {

		DB.Table("vapi_records").Where("id = ?", f.ID).Updates(VapiRecord{Status: "deleted", DataRead: data.GetCurrentDate()})

	} else {
		DB.Table("vapi_records").Where("id = ?", f.ID).Updates(VapiRecord{DataRead: data.GetCurrentDate()})

	}
	fmc.Printfln("#ybtManageMessage> #bbtMessage[#gbt%s#bbt]#wbt, #bbtChatID[#gbt%s#bbt]", msg.Message, msg.ChatID)

}
