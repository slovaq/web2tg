package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/slovaq/web2tg/web/API"
	DAL "github.com/slovaq/web2tg/web/DAL"
	"github.com/slovaq/web2tg/web/vapi"
)

var DB = DAL.DB
var debug bool

func init() {
	debug = os.Getenv("DEBUG") != ""
}
func chk(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		if debug {
			log.Printf("[ERROR] [%s:%s:%d] %s ", runtime.FuncForPC(pc).Name(), fn, line, err)
		} else {
			log.Printf("[ERROR] [%s:%d] %s ", fn, line, err)
		}
	}
}
func index(writer http.ResponseWriter, _ *http.Request) {
	var records []DAL.Record
	var cities []DAL.City
	var users []DAL.User
	DB.Find(&records) // select * from Records to &record
	DB.Find(&cities)  // select * from Records to &record
	DB.Find(&users)   // select * from Records to &record

	tmpl, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		_, err := writer.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
		return
	}

	var structData = struct {
		Records []DAL.Record
		Users   []DAL.User
		Cities  []DAL.City
	}{Records: records, Users: users, Cities: cities}
	err = tmpl.Execute(writer, structData)
	if err != nil {
		_, err := writer.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
	}
}

func staticRouter(w http.ResponseWriter, r *http.Request) {
	file := chi.URLParam(r, "file")
	typeFile := chi.URLParam(r, "type")
	switch typeFile {
	case "css":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "text/css")
	case "js":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "application/javascript")
	case "ttf":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "application/x-font-ttf")
	default:
		log.Printf("Undefined type [%s] of file: [%s]\n", typeFile, file)
	}
	path := "./static/" + typeFile + "/" + file
	//	log.Println(path)
	http.ServeFile(w, r, path)
}
func init() {
	err := DB.AutoMigrate(&DAL.Record{})
	chk(err)
	err = DB.AutoMigrate(&DAL.User{})
	chk(err)
	err = DB.AutoMigrate(&DAL.City{})
	chk(err)
}
func reg(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/reg.html", "templates/base.html")
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
		return
	}
	tmpl.Execute(w, nil)
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		fmt.Println("start")
		login, err := r.Cookie("login")
		if err != nil {
			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		} else {
			fmt.Println(login.Value)
		}
		password, err := r.Cookie("password")
		if err != nil {
			fmt.Println(err.Error())
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				http.Redirect(w, r, "/reg", 301)
				return
			}
		} else {
			fmt.Println(password.Value)
		}
		decodedlogin, err := url.QueryUnescape(login.Value)
		if err != nil {

			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		}
		decodedpassword, err := url.QueryUnescape(password.Value)
		if err != nil {

			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		}
		fmt.Printf("(authMiddleware)> user:%s|password: %s\n", decodedlogin, decodedpassword)
		if decodedlogin == "" || decodedpassword == "" {
			tmpl, err := template.ParseFiles("templates/error.html", "templates/base.html")
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				fmt.Fprintf(w, "error %s", err)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
		_, useErr := DAL.GetUser(decodedlogin, decodedpassword)

		if useErr != nil {
			tmpl, err := template.ParseFiles("templates/error.html", "templates/base.html")
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				fmt.Fprintf(w, "error %s", err)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func profile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/profile.html", "templates/base.html")

	var records []DAL.Record
	var cities []DAL.City
	var users []DAL.User
	DB.Find(&records) // select * from Records to &record
	DB.Find(&cities)  // select * from Records to &record
	DB.Find(&users)   // select * from Records to &record

	var structData = struct {
		Records []DAL.Record
		Users   []DAL.User
		Cities  []DAL.City
	}{Records: records, Users: users, Cities: cities}

	err = tmpl.Execute(w, structData)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
		return
	}

}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/", index)
	r.HandleFunc("/profile", profile)
	r.HandleFunc("/reg", reg)
	r.Route("/auth", func(r chi.Router) {
		r.With(authMiddleware).Route("/", func(r chi.Router) {
			r.Get("/", vapi.GetHandler)
			r.Put("/", vapi.PutHandler)
		})
	})
	r.HandleFunc("/static/{type}/{file}", staticRouter)
	r.Route("/vue", vapi.Router)
	r.Route("/api", API.Router)
	err := http.ListenAndServe(":1111", r)
	if err != nil {
		fmt.Printf("Cant start server with error %s \n Exiting..", err)
	}
}
