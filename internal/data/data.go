package data

import (
	"fmt"
	"time"
)

func Weekday() time.Weekday {

	weekday := time.Now().Weekday()
	//	fmt.Println(weekday)
	return weekday
}

func GetTimestamp(dt, tm string) int64 {
	fulltime := dt + "T" + tm + ":00+03:00"
	t, err := time.Parse(time.RFC3339, fulltime)
	if err != nil {
		fmt.Println(err)
	}
	return t.Unix()

}

//GetCurrentDate return current date in format "2006-01-02"
func GetCurrentDate() string {
	currentTime := time.Now()
	//	fmt.Println("YYYY-MM-DD : ", currentTime.Format("2006-01-02"))
	return currentTime.Format("2006-01-02")
}
