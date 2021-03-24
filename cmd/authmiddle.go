package main

import (
	"net/http"

	"github.com/mallvielfrass/coloredPrint/fmc"
	DAL "github.com/slovaq/web2tg/internal/DAL"
	"github.com/slovaq/web2tg/internal/vapi"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//for k, v := range noCacheHeaders {
		//	w.Header().Set(k, v)
		//}
		login, err := vapi.HandleCookie(r.Cookie("login"))
		if err != nil {
			fmc.Printfln("#rbt(HandleCookie)> Error: #ybt%s", err.Error())
			http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
			return
		}
		password, err := vapi.HandleCookie(r.Cookie("password"))
		if err != nil {
			fmc.Printfln("#rbt(HandleCookie)> Error: #ybt%s", err.Error())
			http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
			return
		}

		fmc.Printf("#gbt(authMiddleware)#bbt> #ybtuser:%s|password: %s\n", login, password)

		_, useErr := DAL.GetUser(login, password)
		if useErr != nil {
			fmc.Printfln("#rbt(middleware.getUser)> Error: #ybt%s", err.Error())
			http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
