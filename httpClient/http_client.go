package httpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jericho-yu/aid/str"
	processBar "github.com/schollz/progressbar/v3"
)

type (
	// HttpClient http客户端
	HttpClient struct {
		Err                error
		requestUrl         string
		requestQueries     map[string]string
		requestMethod      string
		requestBody        []byte
		requestHeaders     map[string][]string
		request            *http.Request
		response           *http.Response
		responseBody       []byte
		responseBodyBuffer *bytes.Buffer
		isReady            bool
		cert               []byte
		transport          *http.Transport
		timeoutSecond      int64
	}
)

var HttpClientApp HttpClient

func (*HttpClient) New(url string) *HttpClient       { return NewHttpClient(url) }
func (*HttpClient) NewGet(url string) *HttpClient    { return NewHttpClientGet(url) }
func (*HttpClient) NewPost(url string) *HttpClient   { return NewHttpClientPost(url) }
func (*HttpClient) NewPut(url string) *HttpClient    { return NewHttpClientPut(url) }
func (*HttpClient) NewDelete(url string) *HttpClient { return NewHttpClientDelete(url) }

// NewHttpClient 实例化：http客户端
//
//go:fix 推荐使用New方法
func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		requestUrl:     url,
		requestQueries: map[string]string{},
		requestHeaders: map[string][]string{
			"Accept":       {},
			"Content-Type": {},
		},
		responseBody:       []byte{},
		responseBodyBuffer: bytes.NewBuffer([]byte{}),
	}
}

// NewHttpClientGet 实例化：http客户端get请求
//
//go:fix 推荐使用NewGet方法
func NewHttpClientGet(url string) *HttpClient {
	return NewHttpClient(url).SetMethod(http.MethodGet)
}

// NewHttpClientPost 实例化：http客户端post请求
//
//go:fix 推荐使用NewPost方法
func NewHttpClientPost(url string) *HttpClient {
	return NewHttpClient(url).SetMethod(http.MethodPost)
}

// NewHttpClientPut 实例化：http客户端put请求
//
//go:fix 推荐使用NewPut方法
func NewHttpClientPut(url string) *HttpClient {
	return NewHttpClient(url).SetMethod(http.MethodPut)
}

// NewHttpClientDelete 实例化：http客户端delete请求
//
//go:fix 推荐使用NewDelete方法
func NewHttpClientDelete(url string) *HttpClient {
	return NewHttpClient(url).SetMethod(http.MethodDelete)
}

// SetCert 设置SSL证书
func (my *HttpClient) SetCert(filename string) *HttpClient {
	var e error

	// 读取证书文件
	if my.cert, e = os.ReadFile(filename); e != nil {
		my.Err = e
	}
	return my
}

// SetUrl 设置请求地址
func (my *HttpClient) SetUrl(url string) *HttpClient {
	my.requestUrl = url
	return my
}

// SetMethod 设置请求方法
func (my *HttpClient) SetMethod(method string) *HttpClient {
	my.requestMethod = method
	return my
}

// AddHeaders 设置请求头
func (my *HttpClient) AddHeaders(headers map[string][]string) *HttpClient {
	my.requestHeaders = headers
	return my
}

// SetQueries 设置请求参数
func (my *HttpClient) SetQueries(queries map[string]string) *HttpClient {
	my.requestQueries = queries
	return my
}

// SetAuthorization 设置认证
func (my *HttpClient) SetAuthorization(username, password, title string) *HttpClient {
	my.requestHeaders["Authorization"] = []string{fmt.Sprintf("%s %s", title, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))}

	return my
}

// SetBody 设置请求体
func (my *HttpClient) SetBody(body []byte) *HttpClient {
	my.requestBody = body

	return my
}

// SetJsonBody 设置json请求体
func (my *HttpClient) SetJsonBody(body any) *HttpClient {
	my.SetHeaderContentType("json")

	my.requestBody, my.Err = json.Marshal(body)
	return my
}

// SetXmlBody 设置xml请求体
func (my *HttpClient) SetXmlBody(body any) *HttpClient {
	my.SetHeaderContentType("xml")

	my.requestBody, my.Err = xml.Marshal(body)

	return my
}

// SetFormBody 设置表单请求体
func (my *HttpClient) SetFormBody(body map[string]string) *HttpClient {
	my.SetHeaderContentType("form")

	params := url.Values{}
	for k, v := range body {
		params.Add(k, v)
	}
	my.requestBody = []byte(params.Encode())

	return my
}

// SetFormDataBody 设置表单数据请求体
func (my *HttpClient) SetFormDataBody(texts map[string]string, files map[string]string) *HttpClient {
	var (
		e      error
		buffer bytes.Buffer
	)

	my.SetHeaderContentType("form-data")

	writer := multipart.NewWriter(&buffer)

	if len(texts) > 0 {
		for k, v := range texts {
			e = writer.WriteField(k, v)
			if e != nil {
				my.Err = e
				return my
			}
		}
	}

	if len(files) > 0 {
		for k, v := range files {
			fileWriter, _ := writer.CreateFormFile("fileField", k)
			file, e := os.Open(v)
			if e != nil {
				my.Err = e
				return my
			}
			_, e = io.Copy(fileWriter, file)
			if e != nil {
				my.Err = e
				return my
			}
			defer func(file *os.File) {
				if err := file.Close(); err != nil {
					fmt.Printf("Failed to close file: %v", err)
				}
			}(file)
		}
	}

	my.requestBody = []byte(writer.FormDataContentType())

	return my
}

// SetPlainBody 设置纯文本请求体
func (my *HttpClient) SetPlainBody(text string) *HttpClient {
	my.SetHeaderContentType("plain")

	my.requestBody = []byte(text)

	return my
}

// SetHtmlBody 设置html请求体
func (my *HttpClient) SetHtmlBody(text string) *HttpClient {
	my.SetHeaderContentType("html")

	my.requestBody = []byte(text)

	return my
}

// SetCssBody 设置Css请求体
func (my *HttpClient) SetCssBody(text string) *HttpClient {
	my.SetHeaderContentType("css")

	my.requestBody = []byte(text)

	return my
}

// SetJavascriptBody 设置Javascript请求体
func (my *HttpClient) SetJavascriptBody(text string) *HttpClient {
	my.SetHeaderContentType("javascript")

	my.requestBody = []byte(text)

	return my
}

// SetSteamBody 设置二进制文件
func (my *HttpClient) SetSteamBody(filename string) *HttpClient {
	var (
		err  error
		file *os.File
	)

	file, err = os.Open(filename)
	if err != nil {
		my.Err = err
		return my
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}(file)

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	// 创建RequestBodyReader用于读取文件内容
	if size > 1*1024*1024 {
		_, my.Err = io.Copy(my.responseBodyBuffer, file)
		if my.Err != nil {
			return my
		}
		my.requestBody = my.responseBodyBuffer.Bytes()
	} else {
		my.requestBody, err = io.ReadAll(file)
		if err != nil {
			my.Err = err
			return my
		}
	}

	// my.request.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	return my
}

// SetHeaderContentType 设置请求头内容类型
func (my *HttpClient) SetHeaderContentType(key ContentType) *HttpClient {
	if val, ok := ContentTypes[key]; ok {
		my.requestHeaders["Content-Type"] = []string{val}
	}

	return my
}

// AppendHeaderContentType 追加请求头内容类型
func (my *HttpClient) AppendHeaderContentType(keys ...ContentType) *HttpClient {
	values := make([]string, len(keys))
	for k, v := range keys {
		if val, ok := ContentTypes[v]; ok {
			values[k] = val
		}
	}

	if len(my.requestHeaders["Content-Type"]) == 0 {
		my.requestHeaders["Content-Type"] = values
	} else {
		my.requestHeaders["Content-Type"] = append(my.requestHeaders["Content-Type"], values...)
	}

	return my
}

// SetHeaderAccept 设置请求头接受内容类型
func (my *HttpClient) SetHeaderAccept(key Accept) *HttpClient {
	if val, ok := Accepts[key]; ok {
		my.requestHeaders["Accept"] = []string{val}
	}

	return my
}

// AppendHeaderAccept 追加请求头接受内容类型
func (my *HttpClient) AppendHeaderAccept(keys ...Accept) *HttpClient {
	values := make([]string, len(keys))
	for k, v := range keys {
		if val, ok := Accepts[v]; ok {
			values[k] = val
		}
	}

	if len(my.requestHeaders["Accept"]) == 0 {
		my.requestHeaders["Accept"] = values
	} else {
		my.requestHeaders["Accept"] = append(my.requestHeaders["Accept"], values...)
	}

	return my
}

// SetTimeoutSecond 设置超时
func (my *HttpClient) SetTimeoutSecond(timeoutSecond int64) *HttpClient {
	my.timeoutSecond = timeoutSecond

	return my
}

// GetResponse 获取响应对象
func (my *HttpClient) GetResponse() *http.Response { return my.response }

// ParseByContentType 根据响应头Content-Type自动解析响应体
func (my *HttpClient) ParseByContentType(target any) *HttpClient {
	switch ContentType(my.GetResponse().Header.Get("Content-Type")) {
	case ContentTypeJson:
		my.GetResponseJsonBody(target)
	case ContentTypeXml:
		my.GetResponseXmlBody(target)
	}

	return my
}

// GetResponseRawBody 获取原始响应体
func (my *HttpClient) GetResponseRawBody() []byte { return my.responseBody }

// GetResponseJsonBody 获取json格式响应体
func (my *HttpClient) GetResponseJsonBody(target any) *HttpClient {
	if e := json.Unmarshal(my.responseBody, &target); e != nil {
		my.Err = e
	}

	return my
}

// GetResponseXmlBody 获取xml格式响应体
func (my *HttpClient) GetResponseXmlBody(target any) *HttpClient {
	if e := xml.Unmarshal(my.responseBody, &target); e != nil {
		my.Err = e
	}
	return my
}

// SaveResponseSteamFile 保存二进制到文件
//
//go:fix 建议使用Download方法
func (my *HttpClient) SaveResponseSteamFile(filename string) *HttpClient {
	// 创建一个新的文件
	file, err := os.Create(filename)
	if err != nil {
		my.Err = err
		return my
	}

	// 将二进制数据写入文件
	_, err = file.Write(my.responseBody)
	if err != nil {
		my.Err = err
		return my
	}

	my.Err = file.Close()

	return my
}

// GetRequest 获取请求
func (my *HttpClient) GetRequest() *http.Request {
	return my.request
}

// GenerateRequest 生成请求对象
func (my *HttpClient) GenerateRequest() *HttpClient {
	var e error

	my.request, e = http.NewRequest(my.requestMethod, my.requestUrl, bytes.NewReader(my.requestBody))
	if e != nil {
		my.Err = fmt.Errorf("生成请求对象失败：%s", e.Error())
		return my
	}

	// 设置请求头
	my.addHeaders()

	// 设置url参数
	my.setQueries()

	// 检查请求对象
	if my.Err = my.check(); my.Err != nil {
		return my
	}

	// 创建一个新的证书池，并将证书添加到池中
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(my.cert) {
		my.Err = errors.New("生成证书失败")
		return my
	}

	// 创建一个新的TLS配置
	tlsConfig := &tls.Config{RootCAs: certPool}

	// 创建一个新的Transport
	my.transport = &http.Transport{TLSClientConfig: tlsConfig}

	my.isReady = true

	return my
}

// beforeSend 发送请求前置动作
func (my *HttpClient) beforeSend() *http.Client {
	if !my.isReady {
		my.GenerateRequest()
		if my.Err != nil {
			return nil
		}
	}

	my.responseBodyBuffer.Reset() // 重置响应体缓存

	// 发送新的请求
	client := &http.Client{Transport: my.transport}

	// 设置超时
	if my.timeoutSecond > 0 {
		client.Timeout = time.Duration(my.timeoutSecond) * time.Second
	}

	return client
}

// Download 下载文件
func (my *HttpClient) Download(filename, processContent string) *HttpClient {
	client := my.beforeSend()
	if my.Err != nil {
		return my
	}

	if my.response, my.Err = client.Do(my.request); my.Err != nil {
		return my
	} else {
		defer my.response.Body.Close()

		f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		defer f.Close()

		if processContent != "" {
			io.Copy(io.MultiWriter(f, processBar.DefaultBytes(my.response.ContentLength, processContent)), my.response.Body)
		} else {
			io.Copy(f, my.response.Body)
		}

		return my
	}
}

// Send 发送请求
func (my *HttpClient) Send() *HttpClient {
	client := my.beforeSend()
	if my.Err != nil {
		return my
	}

	my.response, my.Err = client.Do(my.request)
	if my.Err != nil {
		return my
	}
	defer my.response.Body.Close()

	// 读取新的响应的主体
	if my.response.ContentLength > 1*1024*1024 { // 1MB
		if _, my.Err = io.Copy(my.responseBodyBuffer, my.response.Body); my.Err != nil {
			my.Err = fmt.Errorf("读取响应体失败：%s", my.Err.Error())
			return my
		}
		my.responseBody = my.responseBodyBuffer.Bytes()
	} else {
		my.responseBody, my.Err = io.ReadAll(my.response.Body)
		if my.Err != nil {
			my.Err = fmt.Errorf("读取响应体失败：%s", my.Err.Error())
			return my
		}
	}

	my.isReady = false

	return my
}

// 检查条件是否满足
func (my *HttpClient) check() error {
	if my.requestUrl == "" {
		return errors.New("url不能为空")
	}

	if my.requestMethod == "" {
		my.requestMethod = http.MethodGet
	}

	return nil
}

// 设置url参数
func (my *HttpClient) setQueries() {
	if len(my.requestQueries) > 0 {
		queries := url.Values{}
		for k, v := range my.requestQueries {
			queries.Add(k, v)
		}

		my.requestUrl = str.BufferApp.NewByString(my.requestUrl).String("?", queries.Encode()).ToString()
	}
}

// 设置请求头
func (my *HttpClient) addHeaders() {
	for k, v := range my.requestHeaders {
		my.request.Header[k] = append(my.request.Header[k], v...)
	}
}
