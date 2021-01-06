package DAL

func GetCity(name string) *City {
	var city City
	if result := DB.Where("name = ?", name).First(&city); result.Error != nil {
		return nil
	}
	return &city
}
func CreateOrGetCity(name string, chatID int) *City {
	var city City
	if result := DB.Where("name = ?", name).First(&city); result.Error != nil {
		city = City{
			Name:   name,
			ChatID: chatID,
		}
		DB.Create(&city)
	}

	return &city
}
