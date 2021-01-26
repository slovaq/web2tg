package vapi

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/DAL"
)

var (
	m         sync.Mutex
	d         int
	CheckDate chan bool
	checkInit chan bool
	layout    = "2021-01-18 17:53"
	records   []VapiRecord
)

//Box Message string, Time int64, Token string, URL string, ID int, User string
type Box struct {
	Message string
	Time    int64
	Token   string
	URL     string
	ID      int
	User    string
}

//Boxs []Box
type Boxs []Box

//IntRange struct {min, max int}
type IntRange struct {
	min, max int
}

//NextRandom get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func (upd *UpdateStorage) dBCheck() {
	fmc.Printfln("#gbtDBCheck")
	m.Lock()
	var user []ClientConfig
	DB.Where("").Find(&user)
	boxT := Boxs{}
	var posts []VapiRecord
	DB.Where("status = 'created'").Find(&posts)
	for v := 0; v < len(posts); v++ {
		for d := 0; d < len(user); d++ {
			if user[d].Login == posts[v].User {
				fulltime := posts[v].Date + "T" + posts[v].Time + ":00+03:00"
				t, err := time.Parse(time.RFC3339, fulltime)
				if err != nil {
					fmt.Println(err)
				}
				timestamp := t.Unix()
				boxT.appendToItself(posts[v].Message, timestamp, user[d].BotToken, user[d].ChatLink, posts[v].ID, user[d].Login)
			}
		}
	}

	*&upd.Box = boxT
	m.Unlock()

}

type UpdateStorage struct {
	UpdateRecord chan bool
	UpdateConfig chan string
	ReadRecord   chan bool
	ReadConfig   chan string
	CheckInit    chan bool
	Updatetoken  chan bool
	Box          []Box
}

func InitChannel(UpdateRecord chan bool, UpdateConfig chan string, ReadRecord chan bool, ReadConfig chan string, CheckInit chan bool, Updatetoken chan bool, Box Boxs) *UpdateStorage {
	return &UpdateStorage{
		UpdateRecord: UpdateRecord,
		UpdateConfig: UpdateConfig,
		ReadRecord:   ReadRecord,
		ReadConfig:   ReadConfig,
		CheckInit:    CheckInit,
		Updatetoken:  Updatetoken,
		Box:          Box,
	}
}
func (upd *UpdateStorage) checkDateCounter() {
	for {
		//fmc.Printfln("#gbtcheckDateCounter")
		time.Sleep(time.Duration(1) * time.Second)
		select {
		case <-upd.ReadRecord:
			fmc.Printfln("#gbtcheck date> #rbtDBcheck")
			upd.dBCheck()
			//	upd.UpdateRecord <- true
		}
	}
}
func (box *Boxs) add(item int64) {
	*box = append(*box, Box{Time: item})
}
func (upd *UpdateStorage) ManageMessage(f Box) {
	unixTimeUTC := time.Unix(f.Time, 0)
	unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
	fmc.Printfln("#rbt read> #bbt Time: #gbt %s", unitTimeInRFC3339)
	//var posts []VapiRecord
	//	DB.Where("status = 'created'").Find()
	msg := MessageTG{
		Message: f.Message,
		ChatID:  f.URL,
	}
	MessageTGChannel <- msg
	DB.Table("vapi_records").Where("id = ?", f.ID).Updates(VapiRecord{Status: "deleted"})

}
func (upd *UpdateStorage) read() {
	for {
		select {
		case <-upd.UpdateRecord:
			m.Lock()
			//sort.Sort(box)
			//fmt.Println("boxlen: ", upd.Box)
			bx := append(Boxs{}, upd.Box...)
			sort.Sort(bx)
			//	bx := append(Boxs{}, *box...)
			upd.ManageMessage(bx[0])
			if len(bx) == 1 {
				upd.Box = Boxs{}
			} else {
				if 1 < len(bx) {
					upd.Box = bx[1:]
				}
			}
			m.Unlock()
		}
	}
}

func (box Boxs) Len() int {
	return len(box)
}

func (box Boxs) Less(i, j int) bool {
	return box[i].Time < box[j].Time
	//return false
}

func (box Boxs) Swap(i, j int) {

	box[i], box[j] = box[j], box[i]
}

func (upd *UpdateStorage) Check() {
	fmc.Printfln("#rbt Run> #gbtCheck")
	for {
		m.Lock()
		//fmc.Printfln("#rbt check lock")
		//sort.Sort(upd.Box)
		bx := append(Boxs{}, upd.Box...)
		sort.Sort(bx)
		//fmt.Println(len(bx))
		if 0 < len(bx) {
			dt := time.Now().Local().Unix()
			//		fmc.Printfln("#rbt check> #gbtbx[0]: %d, realTime:%d", bx[0].Time, dt)
			if (bx[0].Time - dt) < 0 {
				v := true
				upd.UpdateRecord <- v
			}
		}
		m.Unlock()
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (box *Boxs) appendToItself(message string, time int64, token string, url string, id int, user string) {
	z := Box{message, time, token, url, id, user}
	*box = append(*box, z)
}

//Add (message string, timestamp int64, token string, url string, id int, user string)
func (box *Boxs) Add(message string, timestamp int64, token string, url string, id int, user string) {
	m.Lock()
	box.appendToItself(message, timestamp, token, url, id, user)
	m.Unlock()

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

}

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
