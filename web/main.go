package main

import (
	"fmt"
	"github.com/slovaq/web2tg/web/API"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	DAL "github.com/slovaq/web2tg/web/DAL"
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

func main() {
	//	errs := make([]error, 3, 3)
	//errs[0] = DB.AutoMigrate(&DAL.Record{})
	//errs[1] = DB.AutoMigrate(&DAL.User{})
	//errs[2] = DB.AutoMigrate(&DAL.City{})
	err := DB.AutoMigrate(&DAL.Record{})
	chk(err)
	err = DB.AutoMigrate(&DAL.User{})
	chk(err)
	err = DB.AutoMigrate(&DAL.City{})
	chk(err)
	//	for _, err := range errs {
	//		if err != nil {
	//			fmt.Printf("error in main() with text: %s \n", err.Error())
	//		}
	//	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static", http.FileServer(http.Dir("static")))
	r.Route("/", func(r chi.Router) {
		r.Get("/", index)
	})
	r.Route("/api", func(r chi.Router) {
		r.Get("/user_create", API.UserCreate)
		r.Get("/user_get", API.UserGet)

		r.Get("record_get", API.GetRecord)
		r.Get("record_create", API.CreateRecord)
	})
	err = http.ListenAndServe(":1111", r)
	if err != nil {
		fmt.Printf("Cant start server with error %s \n Exiting..", err)
	}
}
