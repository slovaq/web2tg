package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	DAL "github.com/slovaq/web2tg/web/DAL"
)

func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		fmt.Println("start")
		login, err := r.Cookie("login")
		if err != nil {
			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		} else {
			fmt.Println(login.Value)
		}
		password, err := r.Cookie("password")
		if err != nil {
			fmt.Println(err.Error())
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				http.Redirect(w, r, "/reg", 301)
				return
			}
		} else {
			fmt.Println(password.Value)
		}
		decodedlogin, err := url.QueryUnescape(login.Value)
		if err != nil {

			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		}
		decodedpassword, err := url.QueryUnescape(password.Value)
		if err != nil {

			fmt.Printf("[login] error %s\n", err.Error())

			//fmt.Fprintf(w, "[login] error %s", err)
			http.Redirect(w, r, "/reg", 301)
			return

		}
		fmt.Printf("(authMiddleware)> user:%s|password: %s\n", decodedlogin, decodedpassword)
		if decodedlogin == "" || decodedpassword == "" {
			tmpl, err := template.ParseFiles("templates/error.html", "templates/base.html")
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				fmt.Fprintf(w, "error %s", err)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
		_, useErr := DAL.GetUser(decodedlogin, decodedpassword)

		if useErr != nil {
			tmpl, err := template.ParseFiles("templates/error.html", "templates/base.html")
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Printf("error in index() with text: %s \n", err.Error())
				}
				fmt.Fprintf(w, "error %s", err)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
