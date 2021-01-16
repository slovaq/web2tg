package sheduler

import (
	"fmt"
	"net/http"
)

var apiString = "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

// @rsngmbot in telegram...
func sendMessage(token string, text string, chatID int) {
	_, err := http.Get(fmt.Sprintf(apiString, token, chatID, text))
	if err != nil {
		// ... add handler
	}
}
