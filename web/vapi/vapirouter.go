package vapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	login, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
		http.Redirect(w, r, "/reg", 301)
		return

	} else {
		fmt.Println(login.Value)
	}
	//	pet := chi.URLParam(r, "req")
	fmt.Fprintf(w, "Hello, %s!", login.Value)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("put post: %s", pet)))
}
