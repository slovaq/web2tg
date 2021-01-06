package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	DAL "github.com/slovaq/web2tg/web/DAL"
	"html/template"
	"net/http"
)

var DB = DAL.DB

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
	errs := make([]error, 3, 3)
	errs[0] = DB.AutoMigrate(&DAL.Record{})
	errs[1] = DB.AutoMigrate(&DAL.User{})
	errs[2] = DB.AutoMigrate(&DAL.City{})
	for _, err := range errs {
		if err != nil {
			fmt.Printf("error in main() with text: %s \n", err.Error())
		}
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static", http.FileServer(http.Dir("static")))
	r.Route("/", func(r chi.Router) {
		r.Get("/", index)
	})
	err := http.ListenAndServe(":1111", r)
	if err != nil {
		fmt.Printf("Cant start server with error %s \n Exiting..", err)
	}
}