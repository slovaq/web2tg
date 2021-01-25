package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
)

//RecordCreate (w http.ResponseWriter, r *http.Request)
func (upd *UpdateStorage) RecordCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordCreate")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return
	}

	login := logix.Value
	city := r.FormValue("city")
	message := r.FormValue("message")
	period := r.FormValue("period")
	dateTimePicker := r.FormValue("date")

	dt := strings.Split(dateTimePicker, " ")
	layout1 := "03:04PM"
	layout2 := "15:04"

	date := dt[0]
	if dt[2] == "pm" {
		dt[2] = "PM"
	}
	if dt[2] == "am" {
		dt[2] = "AM"
	}
	posttime := dt[1] + dt[2]
	t, err := time.Parse(layout1, posttime)
	if err != nil {
		fmt.Println(err)
		return
	}
	normalTime := t.Format(layout2)
	conf := VapiRecord{
		User:    login,
		Message: message,
		City:    city,
		Date:    date,
		Time:    normalTime,
		Status:  "created",
		Period:  period,
	}
	if result := DB.Create(&conf); result.Error != nil {
		fmt.Printf("conf with login %s is exists", login)
	}
	fmc.Printfln("#rbtRecordCreate> #bbt%s> #gbt%s", login, dateTimePicker)
	upd.ReadRecord <- true
	json.NewEncoder(w).Encode(conf)
}

//RecordGet (w http.ResponseWriter, r *http.Request)
func RecordGet(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordGet")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return

	}

	fmt.Println(logix.Value)
	login := logix.Value
	var posts []VapiRecord
	DB.Where("status = \"created\" and user=?", login).Find(&posts)
	for i := 0; len(posts) < 0; i++ {
		fmt.Println(posts[i])

	}

	w.Header().Set("Content-Type", "application/json")
	//	log.Println("unsorted:", posts)
	sort.Sort(PostSorter(posts))
	//	log.Println("by axis:", posts)

	json.NewEncoder(w).Encode(posts)
}

//RecordDelete (w http.ResponseWriter, r *http.Request)
func (upd *UpdateStorage) RecordDelete(w http.ResponseWriter, r *http.Request) {
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		http.Redirect(w, r, "/reg", 301)
		return
	}
	fmt.Println(logix.Value)

	id := r.FormValue("id")
	login := logix.Value

	DB.Where("id = ? and user=?", id, login).Delete(&VapiRecord{})

	w.Header().Set("Content-Type", "application/json")
	fmc.Println("#gbt delete ok!")

	upd.ReadRecord <- true
	json.NewEncoder(w).Encode("{\"status\":\"ok\"}")

}
