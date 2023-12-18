package DAL

import "fmt"

func GetCity(name string) *City {
	var city City
	if name == "" { return nil }
	if result := DB.Where("name = ?", name).First(&city); result.Error != nil {
		return nil
	}
	return &city
}
func GetListOfCity() ([]City, error) {
	var cities []City
	if result := DB.Find(&cities); result.Error != nil {
		return nil, fmt.Errorf(result.Error.Error())
	}

	return cities, nil
}
func CreateOrGetCity(name string, chatID int) *City {
	var city City
	if name == "" { return nil }
	if result := DB.Where("name = ?", name).First(&city); result.Error != nil {
		city = City{
			Name:   name,
			ChatID: chatID,
		}
		DB.Create(&city)
	}

	return &city
}
