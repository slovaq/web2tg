package internal

import (
	"github.com/slovaq/web2tg/internal/repository"
	"github.com/slovaq/web2tg/internal/repository/sqlite3"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}
func (app *App) Run() error {
	db := sqlite3.NewDriver()
	err := db.Connect("sqlite3", "test.db")
	if err != nil {
		return err
	}
	defer db.Disconnect()
	rep, err := repository.NewRepository(db)
	if err != nil {
		return err
	}
	_ = rep
	return nil
}
