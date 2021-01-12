package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	log.Println(">vapi router")
	r.Get("/index", vIndex)
	r.Get("/create", create)

}

type Post struct {
	Date string
	Body string
}
type Posts struct {
	Count int
	Data  []Post
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi create")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")
}
func vIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/posts.html")
	fmt.Fprintf(w, "v index")
}
