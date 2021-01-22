package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func CreateConfig(login string, city string, chatLink string, token string) (*ClientConfig, string, error) {
	conf := ClientConfig{
		Login:    login,
		City:     city,
		ChatLink: chatLink,
		BotToken: token,
	}
	fmt.Printf("login> %s\n\tcity> %s\n\tchatLink> %s\n\ttoken> %s\n", login, city, chatLink, token)
	var user []ClientConfig
	DB.Where("city = ? and login = ?", city, login).Find(&user)
	fmt.Println(user)
	if len(user) == 0 {
		if result := DB.Create(&conf); result.Error != nil {
			return nil, "", fmt.Errorf("conf %s with login %s is exists", login, chatLink)
		}
		return &conf, "create", nil
	}
	rconf := ClientConfig{
		Login:    login,
		City:     city,
		ChatLink: chatLink,
		BotToken: token,
	}
	DB.Model(&rconf).Where("city = ? and login = ?", city, login).Updates(rconf)
	return &rconf, "update", nil

}

type CreateConfData struct {
	User   *ClientConfig
	Status string
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi get post")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		http.Redirect(w, r, "/reg", 301)
		return

	} else {
		fmt.Println(logix.Value)
	}
	login := logix.Value
	var user []ClientConfig
	DB.Where("login = ?", login).Find(&user)
	for i := 0; len(user) < 0; i++ {
		fmt.Println(user[i])

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}
func GetConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi get conf ")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		http.Redirect(w, r, "/reg", 301)
		return

	} else {
		fmt.Println(logix.Value)
	}
	login := logix.Value
	var user []ClientConfig
	DB.Where("login = ?", login).Find(&user)
	for i := 0; len(user) < 0; i++ {
		fmt.Println(user[i])

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func CreateConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi create")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return

	} else {
		fmt.Println(logix.Value)
	}
	login := logix.Value
	chatLink := r.FormValue("chatLink")
	token := r.FormValue("token")
	city := r.FormValue("city")
	user, status, _ := CreateConfig(login, city, chatLink, token)
	w.Header().Set("Content-Type", "application/json")
	data := &CreateConfData{
		User:   user,
		Status: status,
	}
	json.NewEncoder(w).Encode(data)
}
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/posts.html")
	fmt.Fprintf(w, "v index")
}

func RecordCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordCreate")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return
	}

	fmt.Println(logix.Value)

	login := logix.Value
	city := r.FormValue("city")
	dateTimePicker := r.FormValue("date")

	fmt.Println(dateTimePicker)
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
	fmt.Println(t.Format(layout1))
	fmt.Println(t.Format(layout2))
	normalTime := t.Format(layout2)
	message := r.FormValue("message")
	period := r.FormValue("period")
	conf := VapiRecord{
		User:    login,
		Message: message,
		City:    city,
		Date:    date,
		Time:    normalTime,
		Status:  "created",
		Period:  period,
	}
	fmt.Printf("login> %s\n\tcity> %s\n\t\n", login, city)
	if result := DB.Create(&conf); result.Error != nil {
		fmt.Printf("conf with login %s is exists", login)
	}
	f := true
	checkDate <- f
	json.NewEncoder(w).Encode(conf)
}
func RecordDelete(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordDelete")
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
	f := true
	checkDate <- f
	json.NewEncoder(w).Encode("{'status':'ok'}")
}

type PostSorter []VapiRecord

func RecordGet(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordGet")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
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
	//	data := &CreateConfData{
	//		User:   user,
	//		Status: status,
	//	}
	log.Println("unsorted:", posts)
	sort.Sort(PostSorter(posts))
	log.Println("by axis:", posts)

	json.NewEncoder(w).Encode(posts)
}
