package vapi

import (
	"net/http"
	"net/url"
)

func HandleCookie(hCookie *http.Cookie, err error) (string, error) {
	if err != nil {
		return "", err
	}
	decodedValue, err := url.QueryUnescape(hCookie.Value)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}
