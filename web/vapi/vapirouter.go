package vapi

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/go-chi/chi"
)

type Data struct {
	User string
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	login, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
		http.Redirect(w, r, "/reg", 301)
		return

	}
	decodedlogin, err := url.QueryUnescape(login.Value)
	if err != nil {

		fmt.Printf("[login] error %s\n", err.Error())

		//fmt.Fprintf(w, "[login] error %s", err)
		http.Redirect(w, r, "/reg", 301)
		return

	}
	obj := Data{
		User: decodedlogin,
	}
	tmpl, err := template.ParseFiles("templates/home.html", "templates/base.html")
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
		return
	}
	tmpl.Execute(w, obj)
	//	pet := chi.URLParam(r, "req")
	//fmt.Fprintf(w, "Hello, %s!", login.Value)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("put post: %s", pet)))
}
