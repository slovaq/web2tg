package API

import (
	"encoding/json"
	"fmt"
	"log"
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
	login := request.FormValue("login")
	name := request.FormValue("name")
	password := request.FormValue("password")
	log.Printf("(UserCreate)> login: %s | name: %s | password: %s\n", login, name, password)
	//vars := []string{
	//	request.FormValue("login"),
	//	request.FormValue("name"),
	//	request.FormValue("password"),
	//}
	// Если имя не указано - заполнить его из поля логина
	if login == "" {
		obj.Error = fmt.Errorf("No login passed")
		obj.Success = false
	}
	if name == "" {
		name = login
	}
	// FIXME: Проверки не работают
	if password == "" {
		obj.Error = fmt.Errorf("No password passed")
		obj.Success = false
	}

	//
	user, err := DAL.CreateUser(
		login,
		name,
		password,
	)

	if err != nil {
		obj.Error = err.Error()
		obj.Success = false
		log.Println(err)
	}
	obj.Result = user
	returnData, err := json.Marshal(obj)
	if err != nil {
		writer.Write([]byte("Server error"))
		return
	}
	writer.Write(returnData)

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
