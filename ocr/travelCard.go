package ocr

import (
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

type TravelCard struct {
	Region      string
	ImageUrl    string
	ImageBase64 string
}

func (tc *TravelCard) TravelCardOCR() (resp *ocr.RecognizeTravelCardOCRResponse, err error) {
	if tc.Region == "" {
		tc.Region = defaultRegion
	}
	client, err := getClient(tc.Region)
	if err != nil {
		log.Println(err)
		return
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ocr.NewRecognizeTravelCardOCRRequest()
	request.ImageBase64 = common.StringPtr(tc.ImageBase64)
	request.ImageUrl = common.StringPtr(tc.ImageUrl)

	// 返回的resp是一个RecognizeTravelCardOCRResponse的实例，与请求对象对应
	resp, err = client.RecognizeTravelCardOCR(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Println("An API error has returned:", err)
		return
	}
	return
}
