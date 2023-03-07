package web

import "net/http"

func init() {
	http.HandleFunc("/ocr/travel/card", TravelCardHandler)
	http.HandleFunc("/ocr/health/code", HealthCodeHandler)
}
