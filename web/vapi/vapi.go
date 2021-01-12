package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	ChatLink string
	BotToken string
}

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

func CreateConfig(login string, chatLink string, token string) (*ClientConfig, error) {
	conf := ClientConfig{
		Login:    login,
		ChatLink: chatLink,
		BotToken: token,
	}
	if result := DB.Create(&conf); result.Error != nil {
		return nil, fmt.Errorf("conf %s with login %s is exists", login, chatLink)
	}

	return &conf, nil
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
	user, _ := CreateConfig(login, chatLink, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/posts.html")
	fmt.Fprintf(w, "v index")
}
