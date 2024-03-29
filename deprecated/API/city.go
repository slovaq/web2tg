package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/slovaq/web2tg/internal/DAL"
)

func CityGetList(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	obj := JsonObject{
		Success: false,
		Result:  nil,
		Error:   nil,
	}
	cities, err := DAL.GetListOfCity()
	if err != nil {
		chk(&obj, writer, err)
		return
	}

	obj.Success = true
	obj.Result = cities
	returnData, _ := json.Marshal(obj)
	writer.Write(returnData)
}
func CityCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	error_sended := false
	obj := JsonObject{
		Success: false,
		Result:  nil,
		Error:   nil,
	}
	chatId, err := strconv.Atoi(request.FormValue("chat_id"))
	if err != nil {
		chk(&obj, writer, err)
		error_sended = true
	}
	city := DAL.CreateOrGetCity(request.FormValue("name"), chatId)
	obj.Success = true
	obj.Result = city
	returnData, err := json.Marshal(obj)
	if !error_sended {
		writer.Write(returnData)

	}
}
func CityGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	errorSended := false
	result := JsonObject{
		Success: false,
		Result:  nil,
		Error:   nil,
	}
	city := DAL.GetCity(strings.TrimSpace(request.FormValue("city_name")))
	if city == nil {
		chk(&result, writer, fmt.Errorf("city not found"))
		errorSended = true
	}
	result.Result = city
	result.Success = true
	returnData, err := json.Marshal(result)
	if err != nil {
		chk(&result, writer, err)
		fmt.Println("error when marshaling")
		fmt.Println(err.Error())
		return
	}
	if !errorSended {

		writer.Write(returnData)
	}

}
