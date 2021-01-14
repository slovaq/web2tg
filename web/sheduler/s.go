package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Post struct {
	Message string
	Time    int64
}

type Posts []Post

func (s Posts) Len() int {
	return len(s)
}

func (s Posts) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

func (s Posts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var wg sync.WaitGroup

func read(message string) {
	fmt.Printf("read message> %s", message)
}
func check(box *Posts) {
	defer wg.Done()
	for {

		//fmt.Println("checker>", box)

		sc := append([]Post{}, *box...)
		//	fmt.Println("checker lenght>", len(sc))
		if len(sc) > 0 {
			currentTime := time.Now()
			//fmt.Println(currentTime.Format(time.RFC3339))

			//fmt.Println(currentTime.Unix())

			dia := sc[0].Time - currentTime.Unix()

			//fmt.Println("unixTimeUTC>", unixTimeUTC)
			if dia > 0 {

			} else {
				unixTimeUTC := time.Unix(sc[0].Time, 0)
				fmt.Printf("check>\n\tcurrent time: %s\n\tcurrent time unix: %d\n\tmessage time: %s\n\tmessage timestamp: %d\n\tdia: %d\n",
					currentTime.Format(time.RFC3339),
					currentTime.Unix(),
					unixTimeUTC,
					sc[0].Time,
					dia)
				//time.Sleep(1 * time.Millisecond)
				//if (sc[0].Time-currentTime.Unix())<

				go read(sc[0].Message)
				if len(sc) == 1 {
					*box = []Post{}
				} else {
					copy(*box, sc[1:])
				}

			}

		} else {

			*box = []Post{}
			//	fmt.Println("cp>", box)

		}

		time.Sleep(1 * time.Second)
	}
}
func (box *Posts) add() {
	defer wg.Done()
	for {
		r := rand.Intn(10)
		box.appendToItself("message "+strconv.Itoa(r), 24)
		sort.Sort(box)
		time.Sleep(5 * time.Second)
	}
}
func (h *Posts) appendToItself(message string, time int64) {
	z := Post{message, time}
	*h = append(*h, z)
}

func (box *Posts) httpAdd(w http.ResponseWriter, r *http.Request) {

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
	//fmt.Println(t.Format(time.RFC3339))
	//fmt.Println(t.Unix())
	timestamp := t.Unix()
	fmt.Printf("add> message: %s |full time: %s |timestamp: %d\n", message, fulltime, timestamp)
	box.appendToItself(message, timestamp)
	sort.Sort(box)
	fmt.Fprintln(w, "ok")

}
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("start")

	box := Posts{}
	wg.Add(1)
	//go box.add()
	time.Sleep(1 * time.Second)
	go check(&box)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/add", box.httpAdd) ///add?message=test&date=2021-01-14&time=13:42:10
	http.ListenAndServe(":3000", r)
	wg.Wait()

}
