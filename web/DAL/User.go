package DAL

import (
	"fmt"
)

func GetUser(login string, password string) (*User, error) {
	var user User
	if result := DB.Where("login = ? and password = ?", login, password).First(&user); result.Error != nil {
		return nil, fmt.Errorf("login %s not found", login)
	}

	return &user, nil
}
func CreateUser(login string, name string, password string) (*User, error) {
	user := User{
		Login:    login,
		Password: password,
		Name:     name,
	}
	if result := DB.Create(&user); result.Error != nil {
		return nil, fmt.Errorf("login %s is exists", login)
	}
	return &user, nil
}
