package vapi

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/DAL"
)

var (
	m         sync.Mutex
	stateM    infoMutex
	d         int
	CheckDate chan bool
	checkInit chan bool
	layout    = "2021-01-18 17:53"
	records   []VapiRecord
)

func (st *infoMutex) MutexInfoLocker(wci string) {
	if st.Locker != wci {
		st.Status = true
		st.Locker = wci
		//		fmc.Printfln("#gbtMutexInfoLocker> #bbt Mutex[%s] was locked", st.MutexName)
	} else {
		fmc.Printfln("#rbtMutexInfoLocker> LockErr #ybtMutex[%s]#wbt:#bbt do not want: %s, have:%s", st.MutexName, wci, st.Locker)
	}
}
func (st *infoMutex) MutexInfoUnlocker(wci string) {
	if st.Locker == wci {
		st.Status = false
		st.Locker = ""
		//		fmc.Printfln("#gbtMutexInfoUnlocker> #bbt Mutex[%s] was unlocked", st.MutexName)
	} else {
		fmc.Printfln("#rbtMutexInfoUnlocker> LockErr #ybtMutex[%s]#wbt:#bbt want: %s, have:%s", st.MutexName, wci, st.Locker)
	}
}
func checkerChannel() {
	memory := ""
	for {
		if stateM.Locker != memory {
			v := ""
			if stateM.Status == true {
				v = "true"
			} else {
				v = "false"
			}
			fmc.Printfln("#rbtchecker> #ybtStatus: #bbt%s, #ybtLocker: #bbt%s", v, stateM.Locker)
			memory = stateM.Locker
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

//Initrc  start sheduler module
func (upd *UpdateStorage) Initrc() {
	stateM.MutexName = "stateM"
	stateM.Status = false
	rand.Seed(time.Now().UnixNano())
	fmc.Printfln("#rbt Run> #gbtInnitrc")
	//box.add(2)
	//	f := make(chan bool)
	go upd.Check()
	go upd.checkDateCounter()
	go upd.read()
	upd.ReadRecord <- true
	go checkerChannel()
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