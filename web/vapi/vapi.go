package vapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	log.Println(">vapi router")
	r.Get("/index", vIndex)

}
func vIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi index")
	http.ServeFile(w, r, "vapi/template/home/index.html")
	fmt.Fprintf(w, "v index")
}
