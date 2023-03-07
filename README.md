# 健康码行程码识别

### 配置config文件

**/config/secret/secrets.yaml**

这里配置了腾讯云的id和key（用于连接腾讯云接口）（我打码了）以下网址可以查看

https://console.cloud.tencent.com/cam/capi

```yaml
OcrSecret:
  SecretId: "AKIDGI5cuvfB************wQnvqX0FF0V"
  SecretKey: "LQG1Y6RZUBd*************HdmhRwEI"
```

**/config/conf/confs.yaml**

```yaml
OcrConf:
  Endpoint: ocr.tencentcloudapi.com
  Region: ap-shanghai
```

### 引入viper工具进行解析yaml

```go
import "github.com/spf13/viper"

type confInfo struct {
	Name string
	Type string
	Path string
}
type config struct {
	viper *viper.Viper
}

var (
	Conf   *config
	Secret *config
)

func init() {
	ci := confInfo{
		Name: "confs",
		Type: "yaml",
		Path: "config/conf",
	}
	si := confInfo{
		Name: "secrets",
		Type: "yaml",
		Path: "config/secret",
	}
	Conf = &config{getConf(ci)}
	Secret = &config{getConf(si)}
}

func getConf(ci confInfo) *viper.Viper {
	v := viper.New()
	v.SetConfigName(ci.Name) // 与yaml文件名一致
	v.SetConfigType(ci.Type)
	v.AddConfigPath(ci.Path)
	v.ReadInConfig()
	return v
}

func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}

```

### 通过上述配置，获取腾讯云的认证和ocr的连接

```go
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

```

### 这里没有使用任何框架，直接和http建立连接，利用postman进行调试上传文件和上传文件路径

**文件路径存储到腾讯云中的存储桶中，设置公开读私有写的权限，就可访问**，路径类似如下的格式：

"https://wangzj-1258937592.cos.ap-shanghai.myqcloud.com/RecognizeHealthCodeOCR2.png"

（后期我会把这个桶删除，该路径就不可用了，望周知，因为这个只是限时免费）

**另外，mainHealth和mainTravel下的两个main文件仅用于验证通过上述image_url是否可以成功进行ocr识别，后期image_url是可以通过前端输入的（postman）。因此这两个文件可以直接忽略。**

#### 配置路由

```go
func init() {
	http.HandleFunc("/ocr/travel/card", TravelCardHandler)
	http.HandleFunc("/ocr/health/code", HealthCodeHandler)
}
```

#### 实现路由中的两个方法

**1. TravelCardHandler**

* 从前端接收数据`image_url`或者`image_file`，一般不会直接输入`image_base64`进来
* 该功能只要有imageUrl或者imageBase64其一即可实现，如果没有imageUrl 则上传image_file 获得imageBase64
* 调用` ioutil.ReadAll(file)`可以读取file中的内容，以字节的形式返回，再通过`base64.StdEncoding.EncodeToString(bytes)`将字节编码成string，作为我们需要的image_base64
* 实例化一个行程卡，调用`tc.TravelCardOCR()`进行行程卡的ocr识别，并返回结果

```go
func TravelCardHandler(w http.ResponseWriter, req *http.Request) {
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
	tc := &ocr.TravelCard{
		ImageUrl:    imageUrl,
		ImageBase64: imageBase64,
	}
	resp, err := tc.TravelCardOCR()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, resp.ToJsonString())

}
```

其中TravelCardOCR方法如下：

* 其中region是腾讯云ocr接口所需要的字段，我这里默认设置了“ap-shanghai”，通过该region来和该ocr api建立连接
* 然后通过调用接口实现请求，并将参数给到request，最后进行对request的识别并返回

```go
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
```

**2. HealthCodeHandler **

该方法和上述方法类似，就不在重复赘述。

```go
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
```

```go
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
```

