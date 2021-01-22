package vapi

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
)

var (
	m         sync.Mutex
	d         int
	checkDate chan bool
	layout    = "2021-01-18 17:53"
	records   []VapiRecord
)

type Box struct {
	Message string
	Time    int64
	Token   string
	Url     string
	Id      int
	User    string
}
type Boxs []Box
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}
func (box *Boxs) DBCheck() {
	var user []ClientConfig
	DB.Where("").Find(&user)
	boxT := Boxs{}
	var posts []VapiRecord
	DB.Where("status = 'created'").Find(&posts)
	for v := 0; v < len(posts); v++ {
		for d := 0; d < len(user); d++ {
			if user[d].Login == posts[v].User {
				//fmt.Printf("user:%s %s %s %s\n", user[d].Login, user[d].City, user[d].BotToken, user[d].ChatLink)
				//	fmt.Printf("posts:%s %s %s %s %d\n", posts[v].User, posts[v].Message, posts[v].Date, posts[v].Time, posts[v].Id)
				fulltime := posts[v].Date + "T" + posts[v].Time + ":00+03:00"
				//	tm := records[0].Date + "T" + records[0].Time + ":00+03:00" // from MST to Moscow time zone
				//	fmt.Println(tm)
				t, err := time.Parse(time.RFC3339, fulltime)
				if err != nil {
					fmt.Println(err)
				}
				timestamp := t.Unix()
				boxT.appendToItself(posts[v].Message, timestamp, user[d].BotToken, user[d].ChatLink, posts[v].Id, user[d].Login)
			}
		}

	}
	m.Lock()
	*box = boxT
	m.Unlock()

}
func (box *Boxs) checkDateCounter() {
	for {
		select {
		case <-checkDate:
			box.DBCheck()
		}
	}
}
func (box *Boxs) add(item int64) {
	*box = append(*box, Box{Time: item})
}
func (box *Boxs) read(f chan bool) {
	for {
		select {
		case <-f:
			m.Lock()
			sort.Sort(box)
			bx := append(Boxs{}, *box...)
			unixTimeUTC := time.Unix(bx[0].Time, 0)
			unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
			fmc.Printfln("#rbt read> #bbt Time: #gbt %s", unitTimeInRFC3339)
			if len(bx) == 1 {
				*box = Boxs{}
			} else {
				if 1 < len(bx) {
					*box = bx[1:]
				}
			}
			m.Unlock()
		}
	}
}
func (box *Boxs) check(f chan bool) {
	for {
		m.Lock()
		//fmc.Printfln("#rbt check lock")
		sort.Sort(box)
		bx := append(Boxs{}, *box...)
		if 0 < len(bx) {
			dt := time.Now().Local().Unix()
			if (bx[0].Time - dt) < 0 {
				v := true
				f <- v

			}
		}
		m.Unlock()
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (box *Boxs) write(f chan bool) {
	for {
		r := rand.New(rand.NewSource(55))
		ir := IntRange{0, 30}
		m.Lock()
		bx := append(Boxs{}, *box...)
		dt := (time.Now().Local().Unix()) + int64(ir.NextRandom(r))
		unixTimeUTC := time.Unix(dt, 0)
		unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
		fmc.Printfln("#rbt write> #gbtdate: #ybt%s", unitTimeInRFC3339)
		bx = append(bx, Box{Time: dt})
		*box = bx
		m.Unlock()
		time.Sleep(time.Duration(5) * time.Second)
	}
}
func (h *Boxs) appendToItself(message string, time int64, token string, url string, id int, user string) {
	z := Box{message, time, token, url, id, user}
	*h = append(*h, z)
}

func (box *Boxs) Add(message string, timestamp int64, token string, url string, id int, user string) {
	m.Lock()
	box.appendToItself(message, timestamp, token, url, id, user)
	m.Unlock()

}
func Initrc() {
	rand.Seed(time.Now().UnixNano())
	box := Boxs{}
	//box.add(2)
	f := make(chan bool)
	go box.check(f)
	go box.read(f)

	box.checkDateCounter()
}
