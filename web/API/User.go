package API

import (
	"encoding/json"
	"github.com/slovaq/web2tg/web/DAL"
	"net/http"
)

func UserCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	obj := JsonObject{
		Success: false,
		Result:  nil,
		Error:   nil,
	}
	user, err := DAL.CreateUser(
		request.FormValue("login"),
		request.FormValue("name"),
		request.FormValue("password"),
	)
	if err != nil {
		obj.Error = err.Error()
		obj.Success = false
	}
	obj.Result = user
	return_data, err := json.Marshal(obj)
	if err != nil {
		writer.Write([]byte("Server error"))
	}
	writer.Write(return_data)

}
func UserGet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	obj := JsonObject{
		Success: false,
		Result:  nil,
		Error:   nil,
	}
	user, err := DAL.GetUser(
		request.FormValue("login"),
		request.FormValue("password"),
	)
	if err != nil {
		obj.Error = err.Error()
		obj.Success = false
	}
	obj.Result = user
	return_data, err := json.Marshal(obj)
	if err != nil {
		writer.Write([]byte("Server error"))
	}
	writer.Write(return_data)
}
