package vapi

import (
	"fmt"
	"io/ioutil"
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

	fmt.Printf("Hello, %s!\n", decodedlogin)
	dat, err := ioutil.ReadFile("templates/home.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(dat))
	tmpl, err := template.New("").Delims("[[", "]]").Parse(string(dat))
	//	fmt.Fprintf(w, "[login]  %s", obj)
	tmpl.Execute(w, obj)
}

type Infos struct {
	Title, Content string
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("put post: %s", pet)))
}
