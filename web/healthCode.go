package web

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"recognize_health_travel_code/ocr"
)

func HealthCodeHandler(w http.ResponseWriter, req *http.Request) {
	// 该功能只要有imageUrl或者imageBase64其一即可实现
	// 如果没有imageUrl 则上传image_file 获得imageBase64
	imageUrl := req.FormValue("image_url")
	imageBase64 := req.FormValue("image_base64")
	if imageUrl == "" {
		file, _, err := req.FormFile("image_file")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}
		// 读取file，返回字节
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}
		imageBase64 = base64.StdEncoding.EncodeToString(bytes) // 字节编码成字符串
	}
	hc := &ocr.HealthCode{
		ImageUrl:    imageUrl,
		ImageBase64: imageBase64,
	}
	resp, err := hc.HealthCodeOCR()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, resp.ToJsonString())
}
