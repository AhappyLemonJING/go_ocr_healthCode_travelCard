package ocr

import (
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

type HealthCode struct {
	Region      string
	ImageUrl    string
	ImageBase64 string
}

func (hc *HealthCode) HealthCodeOCR() (resp *ocr.RecognizeHealthCodeOCRResponse, err error) {
	client, err := getClient(hc.Region)
	if err != nil {
		log.Println(err)
		return
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewRecognizeHealthCodeOCRRequest()
	request.ImageBase64 = common.StringPtr(hc.ImageBase64)
	request.ImageUrl = common.StringPtr(hc.ImageUrl)

	// 返回的resp是一个RecognizeHealthCodeOCRResponse的实例，与请求对象对应
	resp, err = client.RecognizeHealthCodeOCR(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("An API error has returned: %s", err)
		return
	}
	return
}
