package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	Date string
	Body string
}
type Posts struct {
	Count int
	Data  []Post
}

func Create(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi create")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")
}
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/posts.html")
	fmt.Fprintf(w, "v index")
}
