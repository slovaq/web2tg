package vapi

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/DAL"
)

//Initrc  start sheduler module
func (upd *UpdateStorage) Initrc() {
	rand.Seed(time.Now().UnixNano())
	fmc.Printfln("#rbt Run> #gbtInnitrc")
	//box.add(2)
	//	f := make(chan bool)
	go upd.Check()
	go upd.checkDateCounter()
	go upd.read()
	upd.ReadRecord <- true

	go upd.initBot()

}
func (upd *UpdateStorage) initBot() {
	fmc.Println("#rbtinitBot")
	var user []DAL.ClientConfig
	DB.Where("").Find(&user)
	initV := 0
	if len(user) != 0 {
		fmt.Println("init>user>", user[0].BotToken)
		fmc.Println("#rbtinitBot>Run bot>")
		go upd.runBot()
		initV = 1
	} else {
		fmc.Println("#rbtinitBot>False Run>")
	}
	for {
		select {
		case <-upd.CheckInit:
			if initV == 0 {
				fmc.Println("#rbtinitBot>Run bot>")
				go upd.runBot()
				initV = 1
			} else {
				fmc.Println("#rbtinitBot>Update Token>")
				upd.Updatetoken <- true
			}

		}
	}
}

//fmt.Println("init>user>", user[0].BotToken)
/*
	if len(user) != 0 {

		fmc.Printfln("user:%s", user[0].BotToken)
		if initV == 0 {
			fmc.Println("#rbtinitBot>Run bot>")
			go runBot()
			initV = 1
		} else {
			//	Updatetoken <- "t"
			for {
				select {
				case <-upd.CheckInit:
					if initV == 0 {
						fmc.Println("#rbtinitBot>Run bot>")
						go runBot()
						initV = 1
					} else {
						fmc.Println("#rbtinitBot>Update Token>")
						Updatetoken <- "t"
					}
				}
			}

		}
	} else {

		for {
			select {
			case <-upd.CheckInit:
				if initV == 0 {
					fmc.Println("#rbtinitBot>Run bot>")
					go runBot()
					initV = 1
				} else {
					fmc.Println("#rbtinitBot>Update Token>")
					Updatetoken <- "t"
				}
			}
		}
	}
*/
