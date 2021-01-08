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
