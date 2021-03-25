package vapi

import (
	"fmt"
	"sort"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/data"
)

func (upd *UpdateStorage) dBCheck() {
	//fmc.Caller()
	fmc.Printfln("#rbtDBCheck> #gbtrun")
	m.Lock()

	//fmc.Printfln("#rbtDBCheck> #gbtm.Lock()")
	var user []ClientConfig
	DB.Where("").Find(&user)
	boxT := Boxs{}
	sql := fmt.Sprintf("%%%s%%", data.Weekday())
	//fmt.Println(sql)
	for d := 0; d < len(user); d++ {
		//fmc.Printfln("#ybt Get Table>\n\t#gbtUser: %s", user[d].Login)
		var posts []VapiRecord

		DB.Where("user=? and( (status = 'created' and period='one') or (status = 'created' and period like ? and  data_read!=?))", user[d].Login, sql, data.GetCurrentDate()).Find(&posts)

		fmc.Printfln("#ybtdBCheck> #gbtlen#bbt(posts#bbt) #wbt= #gbt%d", len(posts))
		for v := 0; v < len(posts); v++ {
			post := posts[v]
			fmc.Printfln("\t\t#ybtPost: #bbtMessage:[#gbt%s#bbt]#wbt, #bbt City[#gbt%s#bbt]#wbt, #bbtDate[#gbt%s %s#bbt]#wbt, #bbtPeriod[#gbt%s#bbt]", post.Message, post.City,
				post.Date, post.Time, post.Period)
			tm := data.GetTimestamp(posts[v].Date, posts[v].Time)
			boxT.appendToItself(posts[v].Message, tm, user[d].BotToken, user[d].ChatLink, posts[v].ID, user[d].Login, post.Status, post.Period, post.Pic)
		}
	}
	//	fmc.Printfln("#rbtDBCheck>#gbt iter closed")
	sort.Sort(boxT)
	upd.Box = boxT
	m.Unlock()
	//fmc.Printfln("#rbtfunc DBCheck> #gbtclosed")
}

func (upd *UpdateStorage) Read() {
	for range upd.UpdateRecord {
		fmc.Printfln("#rbtread> #bbtupd.UpdateRecord")
		upd.dBCheck()
		if 0 < len(upd.Box) {
			if (upd.Box[0].Time - time.Now().Local().Unix()) < 0 {
				m.Lock()
				bx := append(Boxs{}, upd.Box...)
				if 1 <= len(bx) {
					upd.ManageMessage(bx[0])
				}
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
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (box *Boxs) appendToItself(message string, time int64, token string, url string, id int, user string, status string, period string, pic string) {
	z := Box{message, time, token, url, id, user, status, period, pic}
	*box = append(*box, z)
}
