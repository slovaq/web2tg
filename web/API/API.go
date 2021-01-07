package API

import (
	"fmt"
	"github.com/slovaq/web2tg/web/DAL"
	"net/http"
	"strconv"
)

func CityCreate(writer http.ResponseWriter, request *http.Request) {
	chat_id, err := strconv.Atoi(request.FormValue("name"))
	if err != nil {
		writer.Write([]byte("chat_id is not valid"))
		return
	}
	city := DAL.CreateOrGetCity(request.FormValue("name"), chat_id)
	writer.Write([]byte(fmt.Sprint(city)))
}
func CityGet(writer http.ResponseWriter, request *http.Request) {
	city := DAL.GetCity(request.FormValue("city_name"))
	if city == nil {
		writer.Write([]byte("city not found"))
		return
	}
	writer.Write([]byte(fmt.Sprint(city)))

}
func UserCreate(writer http.ResponseWriter, request *http.Request) {
	user, err := DAL.CreateUser(
		request.FormValue("login"),
		request.FormValue("name"),
		request.FormValue("password"),
	)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte(fmt.Sprint(user)))

}
func UserGet(writer http.ResponseWriter, request *http.Request) {
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte(fmt.Sprint(user)))
}
func RecordGet(writer http.ResponseWriter, request *http.Request) {
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		writer.Write([]byte("пшел нахуй)0)\n\n"))
		writer.Write([]byte(err.Error()))
		return
	}
	record_id, err := strconv.Atoi(request.FormValue("record_id"))
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
	record, err := user.GetRecord(record_id)
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
	writer.Write([]byte(fmt.Sprint(record)))
}
func RecordCreate(writer http.ResponseWriter, request *http.Request) {

}

//func GetCity(){}
//func CreateCity(){}