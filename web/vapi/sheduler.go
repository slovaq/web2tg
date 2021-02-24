package vapi

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/web/data"
)

//NextRandom get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func (upd *UpdateStorage) dBCheck() {
	fmc.Caller()
	fmc.Printfln("#rbtDBCheck> #gbtopen")
	m.Lock()

	fmc.Printfln("#rbtDBCheck> #gbtm.Lock()")
	var user []ClientConfig
	DB.Where("").Find(&user)
	boxT := Boxs{}
	sql := fmt.Sprintf("%%%s%%", data.Weekday())
	fmt.Println(sql)
	for d := 0; d < len(user); d++ {
		fmc.Printfln("#ybt Get Table>\n\t#gbtUser: %s", user[d].Login)
		var posts []VapiRecord

		//fmt.Println("sql ", sql)
		DB.Where("user=? and( (status = 'created' and period='one') or (status = 'created' and period like ? and  data_read!=?))", user[d].Login, sql, data.GetCurrentDate()).Find(&posts)
		//	DB.Where("user=? and( (status = 'created' and period='one') )", userstatus[d].Login, sql, data.GetCurrentDate()).Find(&posts)

		fmt.Println("post:", posts)
		for v := 0; v < len(posts); v++ {
			post := posts[v]
			fmc.Printfln("\t\t#ybtPost: #bbtMessage:[#gbt%s#bbt]#wbt, #bbt City[#gbt%s#bbt]#wbt, #bbtDate[#gbt%s %s#bbt]#wbt, #bbtPeriod[#gbt%s#bbt]", post.Message, post.City,
				post.Date, post.Time, post.Period)
			tm := data.GetTimestamp(posts[v].Date, posts[v].Time)
			boxT.appendToItself(posts[v].Message, tm, user[d].BotToken, user[d].ChatLink, posts[v].ID, user[d].Login, post.Status, post.Period, post.Pic)
		}
	}
	//var posts []VapiRecord

	//DB.Where("status = 'created' ").Find(&posts)
	//fmc.Printfln("#gbtDBCheck> posts: %v", posts)
	//for v := 0; v < len(posts); v++ {
	//	fmc.Printfln("#gbtDBCheck> iter post: %d", v)
	//	for d := 0; d < len(user); d++ {
	//		if user[d].Login == posts[v].User {
	//			tm := getTimestamp(posts[v].Date, posts[v].Time)
	//			boxT.appendToItself(posts[v].Message, tm, user[d].BotToken, user[d].ChatLink, posts[v].ID, user[d].Login)
	//		}
	//	}
	//	}status
	fmc.Printfln("#rbtDBCheck>#gbt iter closed")
	sort.Sort(boxT)
	*&upd.Box = boxT
	m.Unlock()
	fmc.Printfln("#rbtfunc DBCheck> #gbtclosed")
}

func InitChannel(UpdateRecord chan bool, UpdateConfig chan string, ReadRecord chan bool, ReadConfig chan string, CheckInit chan bool, Updatetoken chan bool, Box Boxs, msg chan MessageTG) *UpdateStorage {
	return &UpdateStorage{
		UpdateRecord: UpdateRecord,
		UpdateConfig: UpdateConfig,
		ReadRecord:   ReadRecord,
		ReadConfig:   ReadConfig,
		CheckInit:    CheckInit,
		Updatetoken:  Updatetoken,
		Box:          Box,
		MessageTG:    msg,
	}
}
func (upd *UpdateStorage) checkDateCounter() {
	for {
		fmc.Printfln("#ybtcheckDateCounter> #gbtawait update")
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
	msg := MessageTG{
		Message: f.Message,
		ChatID:  f.URL,
		Pic:     f.Pic,
	}
	upd.MessageTG <- msg
	if f.Period == "one" {

		DB.Table("vapi_records").Where("id = ?", f.ID).Updates(VapiRecord{Status: "deleted", DataRead: data.GetCurrentDate()})

	} else {
		DB.Table("vapi_records").Where("id = ?", f.ID).Updates(VapiRecord{DataRead: data.GetCurrentDate()})

	}
	fmc.Printfln("#ybtManageMessage> #bbtMessage[#gbt%s#bbt]#wbt, #bbtChatID[#gbt%s#bbt]", msg.Message, msg.ChatID)

}

func (upd *UpdateStorage) read() {
	for {
		select {
		case <-upd.UpdateRecord:
			fmc.Printfln("#rbtread> #bbtupd.UpdateRecord")
			m.Lock()
			bx := append(Boxs{}, upd.Box...)
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
}

func (box Boxs) Swap(i, j int) {

	box[i], box[j] = box[j], box[i]
}

func (upd *UpdateStorage) Check() {
	fmc.Printfln("#rbt Run> #gbtCheck")
	for {
		m.Lock()
		if 0 < len(upd.Box) {
			if (upd.Box[0].Time - time.Now().Local().Unix()) < 0 {
				upd.UpdateRecord <- true
			}
		}
		m.Unlock()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func (box *Boxs) appendToItself(message string, time int64, token string, url string, id int, user string, status string, period string, pic string) {
	z := Box{message, time, token, url, id, user, status, period, pic}
	*box = append(*box, z)
}

//Add (message string, timestamp int64, token string, url string, id int, user string)
func (box *Boxs) Add(message string, timestamp int64, token string, url string, id int, user string, status string, period string, pic string) {
	m.Lock()
	box.appendToItself(message, timestamp, token, url, id, user, status, period, pic)
	m.Unlock()

}
