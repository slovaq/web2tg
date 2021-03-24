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
	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/API"
	"github.com/slovaq/web2tg/internal/DAL"
	"github.com/slovaq/web2tg/internal/bot"

	"github.com/slovaq/web2tg/internal/vapi"
)

var debug bool

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
	path := "./web/static/" + typeFile + "/" + file
	//	log.Println(path)
	http.ServeFile(w, r, path)
}
func init() {
	debug = os.Getenv("DEBUG") != ""
	err := DAL.DB.AutoMigrate(&DAL.Record{})
	chk(err)
	err = DAL.DB.AutoMigrate(&DAL.User{})
	chk(err)
	err = DAL.DB.AutoMigrate(&DAL.City{})
	chk(err)
	err = DAL.DB.AutoMigrate(&DAL.ClientConfig{})
	chk(err)
	err = DAL.DB.AutoMigrate(&vapi.VapiRecord{})
	chk(err)
	err = DAL.DB.AutoMigrate(&vapi.ClientConfig{})
	chk(err)
	err = DAL.DB.AutoMigrate(&vapi.Link{})
	chk(err)

}
func registration(w http.ResponseWriter, r *http.Request) {
	//check if user already authorized
	login, err := vapi.HandleCookie(r.Cookie("login"))
	if err != nil {
		fmc.Printfln("#gbt(registration)> Check: #ybt%s", err.Error())
		//http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		//return
	}
	password, err := vapi.HandleCookie(r.Cookie("password"))
	if err != nil {
		fmc.Printfln("#gbt(registration)> Check: #ybt%s", err.Error())
		//http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		//return
	}
	_, useErr := DAL.GetUser(login, password)
	if useErr != nil {
		fmc.Printfln("#gbt(registration)> Check: #ybtUser not authorized")
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return
	} else {
		fmc.Printfln("#gbt(registration)> Check: Error: #ybtUser #bbt[#gbt%s#bbt]#ybt already authorized", login)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	//____________________________________________
	tmpl, err := template.ParseFiles("web/templates/reg.html", "web/templates/base.html")
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("error in index() with text: %s \n", err.Error())
		}
		return
	}
	tmpl.Execute(w, nil)
}

type Box struct {
	Message string
	Time    int64
	Token   string
	URL     string
	ID      int
	User    string
}

func about(writer http.ResponseWriter, request *http.Request) {
	envVariables := func() string {
		x := "<br>Environment Variables"
		for _, e := range os.Environ() {
			x += "<br>"
			x += e
		}
		return x
	}()
	_, _ = writer.Write([]byte("" +
		"<html><head><title>Hello, my friend!</title></head><body>" +
		"<img src=/static/ttf/egg.png>" +
		envVariables +
		"</body></html>"))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth", http.StatusMovedPermanently)
}
func main() {
	r := chi.NewRouter()
	UpdateRecord := make(chan bool)
	UpdateConfig := make(chan string)
	ReadRecord := make(chan bool)
	UpdateToken := make(chan bool)
	ReadConfig := make(chan string)
	checkInit := make(chan bool)
	msg := make(chan bot.MessageTG)
	box := vapi.Boxs{}
	upd := vapi.InitChannel(UpdateRecord, UpdateConfig, ReadRecord, ReadConfig, checkInit, UpdateToken, box, msg)

	go upd.Initrc()

	r.Use(middleware.Logger)
	//	go sheduler.Listen()
	r.HandleFunc("/", index) //redirect to auth/not auth methods

	//only to this method must be acces for not authorized user
	r.HandleFunc("/reg", registration)
	r.HandleFunc("/user_create", API.UserCreate)
	r.HandleFunc("/static/{type}/{file}", staticRouter)
	r.HandleFunc("/about", about)

	//r.Route("/api", API.Router)
	//to this methods must be acces for only authorized user
	r.Route("/auth", func(r chi.Router) {
		r.With(authMiddleware).Route("/", func(r chi.Router) {
			r.Get("/", vapi.GetHandler)
			//r.Post("/", vapi.PutHandler)
			r.Get("/index", vapi.Index)
			r.Get("/create_config", upd.CreateConf) // /auth/create_config?chatLink=@alalgdfgfdga&token=botfathertokenegbcgbcg&city=test
			r.Get("/get_config", vapi.GetConf)
			r.Get("/get_post", vapi.GetPost)

			r.HandleFunc("/record_get", vapi.RecordGet)
			r.HandleFunc("/record_delete", upd.RecordDelete)
			r.HandleFunc("/record_create", upd.RecordCreate)

			r.Get("/user_get", API.UserGet)

			r.Get("/city_create", API.CityCreate)
			r.Get("/city_get", API.CityGet)
			r.Get("/city_getAll", API.CityGetList)

		})
	})

	err := http.ListenAndServe(":1111", r)
	if err != nil {
		fmt.Printf("Cant start server with error %s \n Exiting..", err)
	}
}
