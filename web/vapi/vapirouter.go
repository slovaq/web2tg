package vapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("get get: %s", pet)))
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("put post: %s", pet)))
}
