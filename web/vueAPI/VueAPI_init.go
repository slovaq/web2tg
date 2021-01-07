package vueAPI

import (
	"net/http"
)

func vIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/home/index.html")
}
