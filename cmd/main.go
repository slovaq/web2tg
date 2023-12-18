package main

import "github.com/slovaq/web2tg/internal"

func main() {
	app := internal.NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
