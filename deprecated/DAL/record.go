package DAL

import (
	"fmt"
	"time"
)

func (u *User) GetRecord(recordID int) (*Record, error) {
	var record Record
	if result := DB.First(&record, recordID); result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}
func (u *User) CreateRecord(date time.Time, dest []City, text string) (*Record, error) {
	record := Record{
		Date:        date,
		Destination: dest,
		Text:        text,
	}
	if result := DB.Create(&record); result.Error != nil {
		return nil, fmt.Errorf("Cant create record ::  %s ", result.Error.Error())
	}
	return &record, nil
}
