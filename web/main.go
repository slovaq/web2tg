package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

var DB, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

func index(writer http.ResponseWriter, request *http.Request) {
	var records []Record
	var cities []City
	var users []User
	DB.Find(&records) // select * from Records to &record
	DB.Find(&cities)  // select * from Records to &record
	DB.Find(&users)   // select * from Records to &record

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	structData := struct {
		Records []Record
		Users   []User
		Cities  []City
	}{Records: records, Users: users, Cities: cities}
	err = tmpl.Execute(writer, structData)
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
}
func main() {
	DB.AutoMigrate(&Record{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&City{})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static", http.FileServer(http.Dir("static")))
	r.Route("/", func(r chi.Router) {
		r.Get("/", index)
	})
	http.ListenAndServe(":1111", r)
}
