package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//GetPost (w http.ResponseWriter, r *http.Request)
func GetPost(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi get post")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return
	}
	fmt.Println(logix.Value)

	login := logix.Value
	var user []ClientConfig
	DB.Where("login = ?", login).Find(&user)
	for i, x := range user {
		fmt.Printf("%d: %v", i, x)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

//Index (w http.ResponseWriter, r *http.Request)
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/posts.html")
	//fmt.Fprintf(w, "v index")
}
