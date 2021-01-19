package vapi

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/slovaq/web2tg/web/DAL"

	"github.com/fatih/color"
)

type PostS struct {
	Message string
	Time    int64
	Token   string
	Url     string
	Id      int
	User    string
}

type PostsS []PostS

func (s PostsS) Len() int {
	return len(s)
}

func (s PostsS) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

func (s PostsS) Swap(i, j int) {

	s[i], s[j] = s[j], s[i]
}

var wg sync.WaitGroup

func read(post PostS) {
	d := color.New(color.FgCyan, color.Bold)
	y := color.New(color.FgYellow, color.Bold)
	y.Printf("read message> ")
	d.Printf("%s %d %s %s\n", post.Message, post.Id, post.Token, post.Url)
	sendMessage(post.Token, post.Url, post.Message)
	DB.Where("id = ? and user=?", post.Id, post.User).Delete(&VapiRecord{})

}
func check(box *PostsS) {
	defer wg.Done()

	for {
		//fmt.Println("checker>", box)
		for {
			if block == true {
				block = false

				//fmt.Printf("block>check >false\n")
				sort.Sort(box)
				//fmt.Println("box>", box)
				sc := append([]PostS{}, *box...)
				//fmt.Println("checker lenght>", len(sc))
				if len(sc) > 0 {
					fmt.Println("checker lenght if len(sc) > 0>", len(sc))
					currentTime := time.Now()
					//fmt.Println(currentTime.Format(time.RFC3339))

					//fmt.Println(currentTime.Unix())

					dia := sc[0].Time - currentTime.Unix()

					fmt.Println("unixTimeUTC>")
					if dia > 0 {
						fmt.Println("dia>", dia)
					} else {
						unixTimeUTC := time.Unix(sc[0].Time, 0)
						fmt.Printf("check>\n\tcurrent time: %s\n\tcurrent time unix: %d\n\tmessage time: %s\n\tmessage timestamp: %d\n\tdia: %d\n",
							currentTime.Format(time.RFC3339),
							currentTime.Unix(),
							unixTimeUTC,
							sc[0].Time,
							dia)

						//go sendMessage(os.Getenv("TOKEN"), sc[0].Message, -404429873)
						go read(sc[0])
						if len(sc) == 1 {
							time.Sleep(10 * time.Millisecond)
							*box = []PostS{}
						} else {
							time.Sleep(1000 * time.Millisecond)
							fmt.Println("sc>", sc)
							//copy(*box, sc[1:])
							*box = sc[1:]
							fmt.Println("box>", box)
						}

					}

				} else {

					//*box = []PostS{}
					//	fmt.Println("cp>", box)

				}

				block = true
				//		fmt.Printf("block>check>true\n")
				time.Sleep(100 * time.Millisecond)
				//		fmt.Println(">break")
				break

			} else {
				//fmt.Printf("block>check>else>true\n")
				time.Sleep(20 * time.Millisecond)
			}
		}

	}
}
func (box *PostsS) add() {
	defer wg.Done()
	for {
		r := rand.Intn(10)
		box.appendToItself("message "+strconv.Itoa(r), 24, "", "", 0, "")

		time.Sleep(5 * time.Second)
	}
}
func (h *PostsS) appendToItself(message string, time int64, token string, url string, id int, user string) {
	z := PostS{message, time, token, url, id, user}
	*h = append(*h, z)
}

var (
	layout                   = "2021-01-18 17:53"
	dateWhenSelectedLastTime time.Time
	records                  []VapiRecord
	block                    bool
)

func (box *PostsS) NohttpAdd(message string, timestamp int64, token string, url string, id int, user string) {
	for {
		if block == true {
			block = false
			//fmt.Printf("block>httpAdd>false\n")
			box.appendToItself(message, timestamp, token, url, id, user)
			block = true
			//	fmt.Printf("block>httpAdd>true\n")
			time.Sleep(100 * time.Millisecond)
			break
		} else {
			//	fmt.Printf("block>httpAdd>sleep\n")
			r := rand.Intn(10)
			time.Sleep(time.Duration(2+r) * time.Millisecond)
		}
	}

}
func (box *PostsS) httpAdd(w http.ResponseWriter, r *http.Request) {

	message := r.FormValue("message")

	date := r.FormValue("date")
	timeValue := r.FormValue("time")

	fulltime := date + "T" + timeValue + "+03:00"
	//fmt.Printf("fulltime> %s\n", fulltime)
	//2021-01-14T13:47:10+03:00
	t, err := time.Parse(time.RFC3339, fulltime)
	if err != nil {
		fmt.Println(err)
	}
	timestamp := t.Unix()
	fmt.Printf("add> message: %s |full time: %s |timestamp: %d\n", message, fulltime, timestamp)
	box.NohttpAdd(message, timestamp, " _", " _", 0, "")
	//sort.Sort(box)
	fmt.Fprintln(w, "ok")

}

func (box *PostsS) chkDB() {
	for {
		box.checkDB()
		time.Sleep(10 * time.Minute)

	}

}
func (box *PostsS) checkDB() {

	var user []ClientConfig
	DB.Where("").Find(&user)
	boxT := PostsS{}
	for {
		if block == true {
			block = false

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
			*box = boxT

			//copy(*box, boxT)
			fmt.Println("boxT>", boxT)
			fmt.Println("box copy>", box)
			//box = boxT
			block = true
			//	fmt.Printf("block>httpAdd>true\n")
			time.Sleep(10 * time.Millisecond)
			break
		} else {
			//	fmt.Printf("block>httpAdd>sleep\n")
			r := rand.Intn(10)
			time.Sleep(time.Duration(2+r) * time.Millisecond)
		}
	}

	//fmt.Println("user:en:", len(user))
	//for v := 0; v < len(user); v++ {
	//	fmt.Printf("user:%s %s %s %s\n", user[v].Login, user[v].City, user[v].BotToken, user[v].ChatLink)
	//
	//}

}

//var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

var z int

func Val() {
	z = 1
	fmt.Println("val> z=1")
}

func InitX() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("start")
	go checkBots()
	//	RFC3339local := "2021-01-14T13:47:10+03:00"
	block = true
	DAL.DB.Take(&records)
	if len(records) != 0 {

		dateWhenSelectedLastTime = time.Now()
		fmt.Printf("records[0].Date> %s|records[0].Time > %s\n ", records[0].Date, records[0].Time)
		fulltime := records[0].Date + "T" + records[0].Time + ":00+03:00"
		//	tm := records[0].Date + "T" + records[0].Time + ":00+03:00" // from MST to Moscow time zone
		//	fmt.Println(tm)
		t, err := time.Parse(time.RFC3339, fulltime)
		if err != nil {
			fmt.Println(err)
		}
		timestamp := t.Unix()
		fmt.Printf("Time: %d-%02d-%02d %02d:%02d:%02d\n\ttimestamp>%d\n",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), timestamp)
		fmt.Printf("message:%s timestamp:%d\n", records[0].Message, timestamp)
	}

	//go box.chkDB()
	box := PostsS{}
	wg.Add(1)
	go check(&box)
	//box.NohttpAdd(fmt.Sprintf("Работаем V!"), timestamp, "971062635:AAGez0zTvZbOlxIqHHVrGF5gnX0mp17nfKs", "https://t.me/joinchat/GBscMQKZaDWZNqBZ",3)

	z = 0
	for {
		if z != 0 {
			fmt.Println("z != 0 ")
			box.checkDB()
			z = 0
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	//box.NohttpAdd(records[0].Message, timestamp)

	//r := chi.NewRouter()
	//	r.Use(middleware.Logger)

	//	r.HandleFunc("/add", box.httpAdd) ///add?message=test&date=2021-01-14&time=13:42:10
	//	http.ListenAndServe(":3001", r)
	//	wg.Wait()

}
