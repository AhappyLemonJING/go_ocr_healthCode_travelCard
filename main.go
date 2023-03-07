package main

import (
	"net/http"
	_ "recognize_health_travel_code/web"
)

func main() {
	http.ListenAndServe(":80", nil)
}
