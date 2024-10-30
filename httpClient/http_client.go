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
)

type (
	// HttpClient http客户端
	HttpClient struct {
		Err            error
		requestUrl     string
		requestQueries map[string]string
		requestMethod  string
		requestBody    []byte
		requestHeaders map[string][]string
		request        *http.Request
		response       *http.Response
		responseBody   []byte
		isReady        bool
		cert           []byte
		transport      *http.Transport
		timeoutSecond  int64
	}
)

// New 实例化：http客户端
func New(url string) *HttpClient {
	return &HttpClient{
		requestUrl:     url,
		requestQueries: map[string]string{},
		requestHeaders: map[string][]string{},
		responseBody:   []byte{},
	}
}

// NewGet 实例化：http客户端get请求
func NewGet(url string) *HttpClient {
	return New(url).SetMethod(http.MethodGet)
}

// NewPost 实例化：http客户端post请求
func NewPost(url string) *HttpClient {
	return New(url).SetMethod(http.MethodPost)
}

// NewPut 实例化：http客户端put请求
func NewPut(url string) *HttpClient {
	return New(url).SetMethod(http.MethodPut)
}

// NewDelete 实例化：http客户端delete请求
func NewDelete(url string) *HttpClient {
	return New(url).SetMethod(http.MethodDelete)
}

// SetCert 设置SSL证书
func (r *HttpClient) SetCert(filename string) *HttpClient {
	var e error

	// 读取证书文件
	if r.cert, e = os.ReadFile(filename); e != nil {
		r.Err = e
	}
	return r
}

// SetUrl 设置请求地址
func (r *HttpClient) SetUrl(url string) *HttpClient {
	r.requestUrl = url
	return r
}

// SetMethod 设置请求方法
func (r *HttpClient) SetMethod(method string) *HttpClient {
	r.requestMethod = method
	return r
}

// AddHeaders 设置请求头
func (r *HttpClient) AddHeaders(headers map[string][]string) *HttpClient {
	r.requestHeaders = headers
	return r
}

// SetQueries 设置请求参数
func (r *HttpClient) SetQueries(queries map[string]string) *HttpClient {
	r.requestQueries = queries
	return r
}

// SetAuthorization 设置认证
func (r *HttpClient) SetAuthorization(username, password, title string) *HttpClient {
	r.requestHeaders["Authorization"] = []string{fmt.Sprintf("%s %s", title, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))}

	return r
}

// SetBody 设置请求体
func (r *HttpClient) SetBody(body []byte) *HttpClient {
	r.requestBody = body
	return r
}

// SetJsonBody 设置json请求体
func (r *HttpClient) SetJsonBody(body any) *HttpClient {
	r.SetHeaderContentType("json")

	r.requestBody, r.Err = json.Marshal(body)
	return r
}

// SetXmlBody 设置xml请求体
func (r *HttpClient) SetXmlBody(body any) *HttpClient {
	r.SetHeaderContentType("xml")

	r.requestBody, r.Err = xml.Marshal(body)
	return r
}

// SetFormBody 设置表单请求体
func (r *HttpClient) SetFormBody(body map[string]string) *HttpClient {
	r.SetHeaderContentType("form")

	params := url.Values{}
	for k, v := range body {
		params.Add(k, v)
	}
	r.requestBody = []byte(params.Encode())
	return r
}

// SetFormDataBody 设置表单数据请求体
func (r *HttpClient) SetFormDataBody(texts map[string]string, files map[string]string) *HttpClient {
	var (
		e      error
		buffer bytes.Buffer
	)

	r.SetHeaderContentType("form-data")

	writer := multipart.NewWriter(&buffer)

	if len(texts) > 0 {
		for k, v := range texts {
			e = writer.WriteField(k, v)
			if e != nil {
				r.Err = e
				return r
			}
		}
	}

	if len(files) > 0 {
		for k, v := range files {
			fileWriter, _ := writer.CreateFormFile("fileField", k)
			file, _ := os.Open(v)
			_, e = io.Copy(fileWriter, file)
			if e != nil {
				r.Err = e
				return r
			}
			defer func(file *os.File) {
				e = file.Close()
				if e != nil {
					panic(e)
				}
			}(file)
		}
	}

	r.requestBody = []byte(writer.FormDataContentType())

	return r
}

// SetPlainBody 设置纯文本请求体
func (r *HttpClient) SetPlainBody(text string) *HttpClient {
	r.SetHeaderContentType("plain")

	r.requestBody = []byte(text)

	return r
}

// SetHtmlBody 设置html请求体
func (r *HttpClient) SetHtmlBody(text string) *HttpClient {
	r.SetHeaderContentType("html")

	r.requestBody = []byte(text)

	return r
}

// SetCssBody 设置Css请求体
func (r *HttpClient) SetCssBody(text string) *HttpClient {
	r.SetHeaderContentType("css")

	r.requestBody = []byte(text)

	return r
}

// SetJavascriptBody 设置Javascript请求体
func (r *HttpClient) SetJavascriptBody(text string) *HttpClient {
	r.SetHeaderContentType("javascript")

	r.requestBody = []byte(text)

	return r
}

// SetSteamBody 设置二进制文件
func (r *HttpClient) SetSteamBody(filename string) *HttpClient {
	var (
		err  error
		file *os.File
	)

	file, err = os.Open(filename)
	if err != nil {
		r.Err = err
		return r
	}
	defer file.Close()

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	// 创建RequestBodyReader用于读取文件内容
	r.requestBody, err = io.ReadAll(file)
	if err != nil {
		r.Err = err
		return r
	}
	r.request.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	return r
	// r.SetHeaderContentType("steam")

	// fileData, err = os.ReadFile(filename)
	// if err != nil {
	// 	r.Err = err
	// 	return r
	// }
	// r.requestBody = fileData

	// return r
}

// SetHeaderContentType 设置请求头内容类型
func (r *HttpClient) SetHeaderContentType(key string) *HttpClient {
	value := ContentType{}.GetValue(key)
	if value != "" {
		r.requestHeaders["Content-Type"] = []string{value}
	}

	return r
}

// SetHeaderAccept 设置请求头接受内容类型
func (r *HttpClient) SetHeaderAccept(key string) *HttpClient {
	value := Accept{}.GetValue(key)
	if value != "" {
		r.requestHeaders["Accept"] = []string{value}
	}

	return r
}

// SetTimeoutSecond 设置超时
func (r *HttpClient) SetTimeoutSecond(timeoutSecond int64) *HttpClient {
	r.timeoutSecond = timeoutSecond

	return r
}

// GetResponse 获取响应对象
func (r *HttpClient) GetResponse() *http.Response {
	return r.response
}

// ParseByContentType 根据响应头Content-Type自动解析响应体
func (r *HttpClient) ParseByContentType(target any) *HttpClient {
	switch r.GetResponse().Header.Get("Content-Type") {
	case ContentTypeJson:
		r.GetResponseJsonBody(target)
	case ContentTypeXml:
		r.GetResponseXmlBody(target)
	}
	return r
}

// GetResponseRawBody 获取原始响应体
func (r *HttpClient) GetResponseRawBody() []byte {
	return r.responseBody
}

// GetResponseJsonBody 获取json格式响应体
func (r *HttpClient) GetResponseJsonBody(target any) *HttpClient {
	if e := json.Unmarshal(r.responseBody, &target); e != nil {
		r.Err = e
	}
	return r
}

// GetResponseXmlBody 获取xml格式响应体
func (r *HttpClient) GetResponseXmlBody(target any) *HttpClient {
	if e := xml.Unmarshal(r.responseBody, &target); e != nil {
		r.Err = e
	}
	return r
}

// SaveResponseSteamFile 保存二进制到文件
func (r *HttpClient) SaveResponseSteamFile(filename string) *HttpClient {
	// 创建一个新的文件
	file, err := os.Create(filename)
	if err != nil {
		r.Err = err
		return r
	}
	defer func() { file.Close() }()

	// 将二进制数据写入文件
	_, err = file.Write(r.responseBody)
	if err != nil {
		r.Err = err
		return r
	}

	return r
}

// GetRequest 获取请求
func (r *HttpClient) GetRequest() *http.Request {
	return r.request
}

// GenerateRequest 生成请求对象
func (r *HttpClient) GenerateRequest() *HttpClient {
	var e error

	r.request, e = http.NewRequest(r.requestMethod, r.requestUrl, bytes.NewReader(r.requestBody))
	if e != nil {
		r.Err = fmt.Errorf("生成请求对象失败：%s", e.Error())
		return r
	}

	// 设置请求头
	r.addHeaders()

	// 设置url参数
	r.setQueries()

	// 检查请求对象
	if r.Err = r.check(); r.Err != nil {
		return r
	}

	// 创建一个新的证书池，并将证书添加到池中
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(r.cert)

	// 创建一个新的TLS配置
	tlsConfig := &tls.Config{RootCAs: certPool}

	// 创建一个新的Transport
	r.transport = &http.Transport{TLSClientConfig: tlsConfig}

	r.isReady = true

	return r
}

// Send 发送请求
func (r *HttpClient) Send() *HttpClient {
	if !r.isReady {
		r.GenerateRequest()
		if r.Err != nil {
			return r
		}
	}

	// 发送新的请求
	client := &http.Client{Transport: r.transport}

	// 设置超时
	if r.timeoutSecond > 0 {
		client.Timeout = time.Duration(r.timeoutSecond) * time.Second
	}

	r.response, r.Err = client.Do(r.request)
	if r.Err != nil {
		r.Err = fmt.Errorf("发送失败：%s", r.Err.Error())
		return r
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.response.Body)

	// 读取新的响应的主体
	r.responseBody, r.Err = io.ReadAll(r.response.Body)
	if r.Err != nil {
		r.Err = fmt.Errorf("读取响应体失败：%s", r.Err.Error())
		return r
	}

	r.isReady = false

	return r
}

// 检查条件是否满足
func (r *HttpClient) check() error {
	if r.requestUrl == "" {
		return errors.New("url不能为空")
	}
	if r.requestMethod == "" {
		r.requestMethod = http.MethodGet
	}
	return nil
}

// 设置url参数
func (r *HttpClient) setQueries() {
	if len(r.requestQueries) > 0 {
		queries := url.Values{}
		for k, v := range r.requestQueries {
			queries.Add(k, v)
		}
		r.requestUrl += "?" + queries.Encode()
	}
}

// 设置请求头
func (r *HttpClient) addHeaders() {
	for k, v := range r.requestHeaders {
		r.request.Header[k] = append(r.request.Header[k], v...)
	}
}
