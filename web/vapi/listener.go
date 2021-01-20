package vapi

import (
	"fmt"
	"net/http"

	"github.com/slovaq/web2tg/web/DAL"
)

var (
	apiString = "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"
	links     []Link
)

// @rsngmbot in telegram...
func sendMessage(token string, url string, text string) {
	DAL.DB.Where("user_link = ?", url).Find(&links)
	fmt.Printf("%s\n\ttoken: %s \n\tlink: %s \n\tmessage> %s\n", redPrint("sendMessage>"), greenPrint(token), yellowPrint(url), bluePrint(text))
	fmt.Println("len links:", len(links))
	fmt.Println("links:", links)

	iurl := fmt.Sprintf(apiString, token, links[0].ChatID, text)
	fmt.Println("iurl>", iurl)
	_, err := http.Get(iurl)
	if err != nil {
		fmt.Printf("err>%s\n", err)
		//	}
	}
}
