package main

import (
	"fmt"
	"recognize_health_travel_code/config"
	"recognize_health_travel_code/ocr"
)

func main() {
	image_url := "https://wangzj-1258937592.cos.ap-shanghai.myqcloud.com/RecognizeHealthCodeOCR2.png"
	healthCode := &ocr.HealthCode{
		ImageUrl: image_url,
	}
	resp, err := healthCode.HealthCodeOCR()
	fmt.Println(resp.ToJsonString(), err)

	fmt.Println(config.Conf.GetString(("OcrConf.Endpoint")))
}
