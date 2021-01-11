package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/slovaq/web2tg/web/DAL"
)

func chk(result *JsonObject, w http.ResponseWriter, err error) {
	// Прим.: если возращает invalid value у RecordGet - у тебя ID неправильный в запросе
	result.Error = err.Error()
	result.Success = false
	return_data, _ := json.Marshal(result)
	w.Write(return_data)
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
	error_sended := false
	result := JsonObject{Success: true, Result: nil, Error: nil}
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		chk(&result, writer, err)
		error_sended = true
	}
	record_id, err := strconv.Atoi(request.FormValue("record_id"))
	if err != nil {
		chk(&result, writer, err)
		error_sended = true

	}
	record, err := user.GetRecord(record_id)
	fmt.Printf("\n%q\n %q \n", record, err)
	if err != nil {
		chk(&result, writer, err)
		error_sended = true

	}
	result.Result = record
	returnData, err := json.Marshal(result)
	if err != nil {
		chk(&result, writer, err)
		fmt.Println("error when marshaling")
		fmt.Println(err.Error())
		return
	}
	if !error_sended {
		writer.Write(returnData)
	}
}
func RecordCreate(writer http.ResponseWriter, request *http.Request) {
	var (
		date   time.Time
		user   *DAL.User
		record *DAL.Record
		cities []DAL.City
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
	citiesDecoded := strings.Split(cities_encoded, ",")
	for _, cityName := range citiesDecoded {
		if city := DAL.GetCity(cityName); city != nil {
			cities = append(cities, *city)
		}

	}
	//layout := "2006-01-02T15:04:05.000Z"

	timer_date := request.FormValue("date")
	text := request.FormValue("text")
	// ждек почини время
	date, err = time.Parse(time.RFC3339, timer_date)
	if err != nil {
		chk(&result, writer, err)
		return
	}
	fmt.Println(date)
	record, err = user.CreateRecord(date, cities, text)
	if err != nil {
		chk(&result, writer, err) // Проверка пользователя на существование
	} else {
		result.Success = true
		result.Result = record
	}

	returnData, err := json.Marshal(result)
	if err != nil {
		chk(&result, writer, err)
		fmt.Println(err.Error())
		return
	}
	writer.Write(returnData)
	/* Получение списка городо в виде
	cities: { "ekb", "moscow } и тд. */
}
