package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	Date string
	Body string
}
type Posts struct {
	Count int
	Data  []Post
}
type ClientConfig struct {
	Login    string
	City     string
	ChatLink string
	BotToken string
}

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

func CreateConfig(login string, city string, chatLink string, token string) (*ClientConfig, string, error) {
	conf := ClientConfig{
		Login:    login,
		City:     city,
		ChatLink: chatLink,
		BotToken: token,
	}
	fmt.Printf("login> %s\n\tcity> %s\n\tchatLink> %s\n\ttoken> %s\n", login, city, chatLink, token)
	//var user ClientConfig
	var user []ClientConfig
	//DB.Raws("select * from client_configs").Scan(&conf)
	//rows, _ := DB.Raw("select * from client_configs where login = ?", conf.Login).Rows()
	//defer rows.Close()
	//	for rows.Next() {

	//		var us ClientConfig
	//		rows.Scan(&us)
	//	user = append(user, us)
	//		fmt.Println(us)
	//		// do something
	//	}
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
	//DB.Model(&ClientConfig{}).Where("login = ?", user[0].Login).Update("name", "hello")
	return &rconf, "update", nil

	//	if result := DB.Create(&conf); result.Error != nil {
	//		return nil, fmt.Errorf("conf %s with login %s is exists", login, chatLink)
	//	}

	//return &conf, nil
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

		//fmt.Fprintf(w, "[login] error %s", err)
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
	//	data := &CreateConfData{
	//		User:   user,
	//		Status: status,
	//	}
	json.NewEncoder(w).Encode(user)
}
func GetConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi get conf ")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
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
	//	data := &CreateConfData{
	//		User:   user,
	//		Status: status,
	//	}
	json.NewEncoder(w).Encode(user)
}
func CreateConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi create")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
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

type VapiRecord struct {
	User    string
	Message string
	City    string
	Date    string
	Id      int
	Time    string
	Status  string
	Period  string
}

func RecordCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordCreate")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return
	} else {
		fmt.Println(logix.Value)
	}
	login := logix.Value
	city := r.FormValue("city")
	//token := r.FormValue("text")
	//	date := r.FormValue("date")
	//time := r.FormValue("time")
	dateTimePicker := r.FormValue("dateTimePicker")
	fmt.Println(dateTimePicker)
	dt := strings.Split(dateTimePicker, " ")
	//form := "2006-01-02T15:04:05"
	date := dt[0]
	time := dt[1]
	//	tm := dt[0] + "T" + dt[1] + "Z" // from MST to Moscow time zone
	//	fmt.Println(tm)
	//	date, err := time.Parse(time.RFC3339, tm)
	//	if err != nil {
	//		fmt.Printf("[login] error %s\n", err.Error())
	//	}
	//	datetime := fmt.Sprintf("Time: %d-%02d-%02d %02d:%02d:%02d-00:00\n",
	//		date.Year(), date.Month(), date.Day(),
	//		date.Hour(), date.Minute(), date.Second())
	//	fmt.Println(datetime)
	//timer_date := request.FormValue("date")
	message := r.FormValue("message")
	period := r.FormValue("period")
	//	data := Record{

	//}
	conf := VapiRecord{
		User:    login,
		Message: message,
		City:    city,
		Date:    date,
		Time:    time,
		Id:      0,
		Period:  period,
		Status:  "created",
	}
	fmt.Printf("login> %s\n\tcity> %s\n\t\n", login, city)
	if result := DB.Create(&conf); result.Error != nil {
		fmt.Errorf("conf with login %s is exists", login)
	}
	json.NewEncoder(w).Encode(conf)
}
func RecordGet(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordGet")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
		http.Redirect(w, r, "/reg", 301)
		return

	} else {
		fmt.Println(logix.Value)
	}
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
	json.NewEncoder(w).Encode(posts)
}
