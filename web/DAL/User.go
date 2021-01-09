package DAL

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)
func SHA256(text string) string {
	algorithm := sha256.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func GetUser(login string, password string) (*User, error) {
	var user User
	hash := SHA256(password)
	if result := DB.Where("login = ? and password = ?", login, hash).First(&user); result.Error != nil {
		return nil, fmt.Errorf("login %s not found", login)
	}
	user.Password = "hidden"
	return &user, nil
}
func CreateUser(login string, name string, password string) (*User, error) {
	user := User{
		Login:    login,
		Password: SHA256(password),
		Name:     name,
	}
	if result := DB.Create(&user); result.Error != nil {
		return nil, fmt.Errorf("login %s is exists", login)
	}
	user.Password = "hidden"
	return &user, nil
}
