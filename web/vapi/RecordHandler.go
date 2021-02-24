package vapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mallvielfrass/coloredPrint/fmc"
)

//RecordCreate (w http.ResponseWriter, r *http.Request)
func (upd *UpdateStorage) RecordCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordCreate")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
		return
	}
	pic := ""
	login := logix.Value
	city := r.FormValue("city")
	message := r.FormValue("message")
	period := r.FormValue("period")
	dateTimePicker := r.FormValue("date")
	//file := r.FormValue("file")
	_, _, fileerr := r.FormFile("file")
	if fileerr != nil {
		fmc.Printfln("#rbtError: #ybt%s", fileerr.Error())
	} else {
		file, fileheader, err := r.FormFile("file")
		if err != nil || fileheader.Size == 0 {
			// file was not sent
		} else {
			// process file

			FileName := "files/" + fileheader.Filename
			f, err := os.Create(FileName)
			if err != nil {
				fmc.Printfln("#rbtError: #ybt%s", err.Error())
			}
			// It's idiomatic to defer a `Close` immediately
			// after opening a file.
			defer f.Close()
			fmt.Println("file: ", file)
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				fmc.Printfln("#rbtError: #ybt%s", err.Error())
			}
			n2, err := f.Write(buf.Bytes())
			if err != nil {
				fmc.Printfln("#rbtError: #ybt%s", err.Error())
			}
			fmt.Printf("wrote %d bytes\n", n2)

			//	io.Copy(f, file)
			pic = FileName
		}
		defer file.Close()
	}

	var links []Link
	DB.Where("").Find(&links)
	//fmt.Println(links[0])
	if len(links) != 0 {
		dt := strings.Split(dateTimePicker, " ")
		layout1 := "03:04PM"
		layout2 := "15:04"

		date := dt[0]
		if dt[2] == "pm" {
			dt[2] = "PM"
		}
		if dt[2] == "am" {
			dt[2] = "AM"
		}
		posttime := dt[1] + dt[2]
		t, err := time.Parse(layout1, posttime)
		if err != nil {
			fmt.Println(err)
			return
		}
		normalTime := t.Format(layout2)
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
		upd.ReadRecord <- true
		json.NewEncoder(w).Encode(conf)

	} else {
		fmc.Println("error:'links not found'")
		f := HandleError{
			HttpError: "link",
		}
		json.NewEncoder(w).Encode(f)
	}

}

type HandleError struct {
	HttpError string
}

//RecordGet (w http.ResponseWriter, r *http.Request)
func RecordGet(w http.ResponseWriter, r *http.Request) {
	log.Println(">vapi RecordGet")
	logix, err := r.Cookie("login")
	if err != nil {
		fmt.Printf("[login] error %s\n", err.Error())
		http.Redirect(w, r, "/reg", 301)
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

		http.Redirect(w, r, "/reg", 301)
		return
	}
	fmt.Println(logix.Value)

	id := r.FormValue("id")
	login := logix.Value

	DB.Where("id = ? and user=?", id, login).Delete(&VapiRecord{})

	w.Header().Set("Content-Type", "application/json")
	fmc.Println("#gbt delete ok!")

	upd.ReadRecord <- true
	json.NewEncoder(w).Encode("{\"status\":\"ok\"}")

}
