package vapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	r.Get("/", vIndex)

}
func vIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/home/index.html")
}
