package API

import (
	"fmt"
	"github.com/slovaq/web2tg/web/DAL"
	"net/http"
)

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
