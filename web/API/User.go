package API

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slovaq/web2tg/web/DAL"
)

func UserCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	obj := JsonObject{
		Success: true,
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
		return
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
	password := request.FormValue("password")
	login := request.FormValue("login")
	if password == "" {
		chk(&obj, writer, fmt.Errorf("no password passed"))
		return
	}
	if login == "" {
		chk(&obj, writer, fmt.Errorf("no login passed"))
		return

	}
	user, err := DAL.GetUser(
		login, password,
	)

	if err != nil {
		chk(&obj, writer, err)
		return
	}
	obj.Result = user
	obj.Success = true
	return_data, err := json.Marshal(obj)
	if err != nil {
		writer.Write([]byte("Server error"))
	}
	writer.Write(return_data)
}
