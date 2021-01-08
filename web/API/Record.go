package API

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/slovaq/web2tg/web/DAL"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func chk(result *JsonObject, w http.ResponseWriter, err error) {
	result.Error = err.Error()
	result.Success = false
	return_data, _ := json.Marshal(result)
	fmt.Println(return_data)
	w.Write([]byte(return_data))
	if err != nil {
		fmt.Println(err.Error())
	}
}

type JsonObject struct {
	Success bool
	Result  interface{}
	Error   interface{}
}

func RecordGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	result := JsonObject{Success: true, Result: nil, Error: nil}
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		chk(&result, writer, err)
	}
	record_id, err := strconv.Atoi(request.FormValue("record_id"))
	if err != nil {
		chk(&result, writer, err)

	}
	record, err := user.GetRecord(record_id)
	if record == nil {
		chk(&result, writer, fmt.Errorf("записи нет"))
	}
	if err != nil {
		chk(&result, writer, err)
	}
	result.Result = record
	return_data, err := json.Marshal(result)
	if err != nil {
		chk(&result, writer, err)
		fmt.Println(err.Error())
		return
	}
	writer.Write([]byte(return_data))
}
func RecordCreate(writer http.ResponseWriter, request *http.Request) {
	var (
		date   time.Time
		user   *DAL.User
		record *DAL.Record
		cities []*DAL.City
	)
	writer.Header().Set("Content-Type", "application/json")
	result := JsonObject{Success: true, Result: nil, Error: nil}
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		chk(&result, writer, err) // Проверка пользователя на существование
	}
	cities_encoded := request.FormValue("cities") // from json comma separated string encoded to []string
	cities_decoded := strings.Split(cities_encoded, ",")
	for _, city_name := range cities_decoded {
		if city := DAL.GetCity(city_name); city != nil {
			cities = append(cities, city)
		}

	}
	//layout := "2006-01-02T15:04:05.000Z"

	//timer_date := request.FormValue("time")
	text := request.FormValue("text")

	/* Получение списка городо в виде
	cities: { "ekb", "moscow } и тд. */
}
