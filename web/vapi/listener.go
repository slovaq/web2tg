package vapi

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/slovaq/web2tg/web/DAL"
)

var apiString = "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

// @rsngmbot in telegram...
func sendMessage(token string, url string, text string) {
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	b := color.New(color.FgBlue, color.Bold).SprintFunc()
	g := color.New(color.FgGreen, color.Bold).SprintFunc()
	var links []Link
	DAL.DB.Where("user_link = ?", url).Find(&links)
	fmt.Printf("%s\n\ttoken: %s \n\tlink: %s \n\tmessage> %s\n", red("sendMessage>"), g(token), yellow(url), b(text))
	iurl := fmt.Sprintf(apiString, token, links[0].ChatID, text)
	fmt.Println("iurl>", iurl)
	_, err := http.Get(iurl)
	if err != nil {
		fmt.Printf("err>%s\n", err)
		//	}
	}
}
