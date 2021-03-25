package vapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mallvielfrass/coloredPrint/fmc"
	"github.com/slovaq/web2tg/internal/data"
)

//RecordCreate (w http.ResponseWriter, r *http.Request)
func prepareImg(file multipart.File, header *multipart.FileHeader) string {
	FileName := "files/" + header.Filename
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
	pic := FileName
	fmt.Println("file: ", header.Size)
	src, err := imaging.Open(pic)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	// Resize srcImage to size = 128x128px using the Lanczos filter.
	g := src.Bounds()

	// Get height and width

	if 4500000 < header.Size || 3061 < g.Dy() || 3061 < g.Dx() {
		height := float64(g.Dy())
		width := float64(g.Dx())

		perc := float64(4999999) / float64(header.Size)
		fmt.Printf(" %f %f \n perc= %f\n", width, height, perc)
		if 4500000 < header.Size {

			height = height * perc
			width = width * perc
		}

		fmt.Printf(" %f %f\n", width, height)

		for {
			if 3060 < height || 3060 < width {
				height = height * 0.9
				width = width * 0.9
			} else {
				break
			}
		}
		h := int(height)
		w := int(width)
		fmt.Printf(" %d %d\n", w, h)
		d := fmt.Sprintf("ffmpeg -y -i %s -vf scale=%d:%d %s", pic, w, h, pic)
		//d := fmt.Sprintf("ffmpeg -y -i %s -vf scale=512:512 %s", pic, pic)
		lsCmd := exec.Command("bash", "-c", d)
		_, err = lsCmd.Output()
		if err != nil {
			panic(err)
		}
	}
	defer file.Close()
	return pic
}
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
	file, fileheader, fileerr := r.FormFile("file")
	if fileerr == nil {
		// process file
		pic = prepareImg(file, fileheader)
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
			upd.ReadRecord <- true
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
				upd.ReadRecord <- true
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
				upd.ReadRecord <- true
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
