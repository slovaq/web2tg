package user

import (
	"github.com/slovaq/web2tg/internal/model"
	"github.com/slovaq/web2tg/internal/repository/sqlite3"
)

type UserSchema struct {
	Login    string `field:"login" index:"true"`
	Password string `field:"password"`
	Name     string `field:"name"`
	Email    string `field:"email"`
}
type User struct {
	driver *sqlite3.ISqlite3
}

func Init(driver *sqlite3.ISqlite3) (*User, error) {
	//create table
	//	fmt.Printf("init user\n")
	tableName := "users"
	err := driver.Migrate(tableName, UserSchema{})
	if err != nil {
		return nil, err
	}

	return &User{
		driver: driver,
	}, nil
}
func (u *User) GetUserByLoginOrName(login string) (*model.User, error) {
	return nil, nil
}
