package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

var db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

func index(writer http.ResponseWriter, request *http.Request) {
	var records []Record
	var cities []City
	var users []User
	db.Find(&records) // select * from records to &record
	db.Find(&cities)  // select * from records to &record
	db.Find(&users)   // select * from records to &record

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	struct_data := struct {
		records []Record
		users   []User
		cities  []City
	}{records: records, users: users, cities: cities}
	tmpl.Execute(writer, struct_data)
}
func main() {
	db.AutoMigrate(&Record{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&City{})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", index)
	})
	http.ListenAndServe(":1111", r)
}
