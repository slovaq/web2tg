package vapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/data"
)

//RecordCreate (w http.ResponseWriter, r *http.Request)

func (upd *UpdateStorage) RecordCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordCreate")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return
	}
	pic := ""
	login := logix.Value
	city := r.FormValue("city")
	message := r.FormValue("message")
	period := r.FormValue("period")
	dateTimePicker := r.FormValue("date")
	//file := r.FormValue("file")
	file, fileheader, fileerr := r.FormFile("file")
	if fileerr == nil {
		// process file
		pic = data.PrepareImg(file, fileheader)
	}

	var links []Link
	DB.Where("").Find(&links)
	if len(links) != 0 {
		dt := strings.Split(dateTimePicker, " ")
		layout1 := "03:04PM"
		layout2 := "15:04"

		date := dt[0]
		dt[2] = strings.ToUpper(dt[2])
		postTime := dt[1] + dt[2]
		t, err := time.Parse(layout1, postTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		normalTime := t.Format(layout2)
		switch period {

		case "one":
			fmt.Println("case: one")
			conf := VapiRecord{
				User:    login,
				Message: message,
				City:    city,
				Date:    date,
				Time:    normalTime,
				Status:  "created",
				Period:  period,
				Pic:     pic,
			}
			if result := DB.Create(&conf); result.Error != nil {
				fmt.Printf("conf with login %s is exists", login)

			}
			fmc.Printfln("#rbtRecordCreate> #bbt%s> #gbt%s", login, dateTimePicker)
			upd.UpdateRecord <- true
			json.NewEncoder(w).Encode(conf)
		case "week":
			checkedNames := r.FormValue("week")
			fmt.Println("ck: ", checkedNames)
			fmt.Println("case: week")
			if strings.Contains(checkedNames, data.Weekday().String()) {
				fmc.Printfln("week contains %s", data.Weekday().String())
				conf := VapiRecord{
					User:    login,
					Message: message,
					City:    city,
					Date:    date,
					Time:    normalTime,
					Status:  "created",
					Period:  checkedNames,
					Pic:     pic,
				}
				if result := DB.Create(&conf); result.Error != nil {
					fmt.Printf("conf with login %s is exists", login)

				}
				fmc.Printfln("#rbtRecordCreate> #bbt%s> #gbt%s", login, dateTimePicker)
				upd.UpdateRecord <- true
				json.NewEncoder(w).Encode(conf)
			} else {

				conf := VapiRecord{
					User:    login,
					Message: message,
					City:    city,
					Date:    date,
					Time:    normalTime,
					Status:  "created",
					Period:  "one",
					Pic:     pic,
				}
				if result := DB.Create(&conf); result.Error != nil {
					fmt.Printf("conf with login %s is exists", login)

				}
				fmc.Printfln("#rbtRecordCreate> #bbt%s> #gbt%s", login, dateTimePicker)
				conf = VapiRecord{
					User:    login,
					Message: message,
					City:    city,
					Date:    date,
					Time:    normalTime,
					Status:  "created",
					Period:  checkedNames,
					Pic:     pic,
				}
				if result := DB.Create(&conf); result.Error != nil {
					fmt.Printf("conf with login %s is exists", login)

				}
				fmc.Printfln("#rbtRecordCreate> #bbt%s> #gbt%s", login, dateTimePicker)
				upd.UpdateRecord <- true
				json.NewEncoder(w).Encode(conf)
			}
		default:
			fmc.Printfln("case: %s not found", period)
		}
		//fmt.Println(data.Weekday())

	} else {
		fmc.Println("error:'links not found'")
		f := HandleError{
			HttpError: "link",
		}
		json.NewEncoder(w).Encode(f)
	}

}

//RecordGet (w http.ResponseWriter, r *http.Request)
func RecordGet(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordGet")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return

	}

	fmt.Println(logix.Value)
	login := logix.Value
	var posts []VapiRecord
	DB.Where("status = \"created\" and user=?", login).Find(&posts)
	for i := 0; len(posts) < 0; i++ {
		fmt.Println(posts[i])

	}

	w.Header().Set("Content-Type", "application/json")
	//	log.Println("unsorted:", posts)
	sort.Sort(PostSorter(posts))
	//	log.Println("by axis:", posts)

	json.NewEncoder(w).Encode(posts)
}

//RecordDelete (w http.ResponseWriter, r *http.Request)
func (upd *UpdateStorage) RecordDelete(w http.ResponseWriter, r *http.Request) {
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())

		http.Redirect(w, r, "/reg", http.StatusMovedPermanently)
		return
	}
	fmt.Println(logix.Value)

	id := r.FormValue("id")
	login := logix.Value

	DB.Where("id = ? and user=?", id, login).Delete(&VapiRecord{})

	w.Header().Set("Content-Type", "application/json")
	fmc.Println("#gbt delete ok!")

	upd.UpdateRecord <- true
	json.NewEncoder(w).Encode("{\"status\":\"ok\"}")

}
