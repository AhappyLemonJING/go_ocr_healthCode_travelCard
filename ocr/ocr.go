package ocr

import (
	"recognize_health_travel_code/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

var credential *common.Credential
var ocrEndPoint string
var defaultRegion string

func init() {
	credential = getCredential()
	ocrEndPoint = config.Conf.GetString("OcrConf.Endpoint")
	defaultRegion = config.Conf.GetString("OcrConf.Region")
}

func getCredential() *common.Credential {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential1 := common.NewCredential(
		config.Secret.GetString("OcrSecret.SecretId"),
		config.Secret.GetString("OcrSecret.SecretKey"),
	)
	return credential1
}

func getClient(region string) (*ocr.Client, error) {
	if region == "" {
		region = defaultRegion
	}
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = string(ocrEndPoint)
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := ocr.NewClient(credential, region, cpf)
	return client, err
}
