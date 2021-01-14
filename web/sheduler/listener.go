package main

import (
	"fmt"
	"net/http"
)

var apiString = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s"

// @rsngmbot in telegram...
func sendMessage(token string, text string, chatID int) {
	resp, err := http.Get(fmt.Sprintf(apiString, token, chatID, text))
	if err != nil {
		// ... add handler
	}
	fmt.Println(resp)
}
