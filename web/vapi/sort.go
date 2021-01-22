package vapi

import (
	"fmt"
	"time"
)

func (s Boxs) Len() int {
	return len(s)
}

func (s Boxs) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

func (s Boxs) Swap(i, j int) {

	s[i], s[j] = s[j], s[i]
}

func (a PostSorter) Len() int      { return len(a) }
func (a PostSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PostSorter) Less(i, j int) bool {
	RFC3339local := "2006-01-02T15:04:05Z"
	aitm := a[i].Date + "T" + a[i].Time + ":00Z" // from MST to Moscow time zone
	aitmdate, err := time.Parse(RFC3339local, aitm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Time: %d-%02d-%02d %02d:%02d:%02d-00:00\n",
		aitmdate.Year(), aitmdate.Month(), aitmdate.Day(),
		aitmdate.Hour(), aitmdate.Minute(), aitmdate.Second())
	fmt.Println(aitmdate.Unix())

	ajtm := a[j].Date + "T" + a[j].Time + ":00Z" // from MST to Moscow time zone
	ajtmdate, err := time.Parse(RFC3339local, ajtm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Time: %d-%02d-%02d %02d:%02d:%02d-00:00\n",
		ajtmdate.Year(), ajtmdate.Month(), ajtmdate.Day(),
		ajtmdate.Hour(), ajtmdate.Minute(), ajtmdate.Second())
	fmt.Println(ajtmdate.Unix())
	return aitmdate.Unix() < ajtmdate.Unix()

}
