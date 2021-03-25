package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//CreateConf (w http.ResponseWriter, r *http.Request)
func (upd *UpdateStorage) CreateConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi create")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return

	}
	fmt.Println(logix.Value)
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
	upd.GobotConnect.CheckInit <- true
	json.NewEncoder(w).Encode(data)
}

//GetConf (w http.ResponseWriter, r *http.Request)
func GetConf(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi get conf ")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return
	}
	fmt.Println(logix.Value)

	login := logix.Value
	var user []ClientConfig
	DB.Where("login = ?", login).Find(&user)
	for i := 0; len(user) < 0; i++ {
		fmt.Println(user[i])

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

//CreateConfig (login string, city string, chatLink string, token string) (*ClientConfig, string, error)
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
