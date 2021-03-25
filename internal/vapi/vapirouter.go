package vapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/mallvielfrass/coloredPrint/fmc"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	login, err := HandleCookie(r.Cookie("login"))
	if err != nil {
		fmc.Printfln("#rbt(HandleCookie)> Error: #ybt%s", err.Error())
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return
	}
	obj := Data{
		User: login,
	}

	fmt.Printf("Hello, %s!\n", login)
	dat, err := ioutil.ReadFile("./web/pages/home/home.html")
	if err != nil {
		fmt.Println(err)
	}
	tmpl, err := template.New("").Delims("[[", "]]").Parse(string(dat))
	if err != nil {
		fmt.Println(err)
	}
	//	fmt.Fprintf(w, "[login]  %s", obj)
	tmpl.Execute(w, obj)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	pet := chi.URLParam(r, "req")
	w.Write([]byte(fmt.Sprintf("put post: %s", pet)))
}
