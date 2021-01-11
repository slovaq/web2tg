package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"

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

	tmpl, err := template.ParseFiles("templates/index.html")
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
	http.ServeFile(w, r, "templates/reg.html")
}
func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("name")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(c.Value)
		}
		if c == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("forbidden!\n")))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/", index)
	r.HandleFunc("/reg", reg)
	r.Route("/auth", func(r chi.Router) {
		r.With(authMiddleware).Route("/{req}", func(r chi.Router) {
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
