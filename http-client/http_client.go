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
)

type (
	// Client http客户端
	Client struct {
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

// NewHttpClient 实例化：http客户端
func NewHttpClient(url string) *Client {
	return &Client{
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
func NewHttpClientGet(url string) *Client {
	return NewHttpClient(url).SetMethod(http.MethodGet)
}

// NewHttpClientPost 实例化：http客户端post请求
func NewHttpClientPost(url string) *Client {
	return NewHttpClient(url).SetMethod(http.MethodPost)
}

// NewHttpClientPut 实例化：http客户端put请求
func NewHttpClientPut(url string) *Client {
	return NewHttpClient(url).SetMethod(http.MethodPut)
}

// NewHttpClientDelete 实例化：http客户端delete请求
func NewHttpClientDelete(url string) *Client {
	return NewHttpClient(url).SetMethod(http.MethodDelete)
}

// SetCert 设置SSL证书
func (my *Client) SetCert(filename string) *Client {
	var e error

	// 读取证书文件
	if my.cert, e = os.ReadFile(filename); e != nil {
		my.Err = e
	}
	return my
}

// SetUrl 设置请求地址
func (my *Client) SetUrl(url string) *Client {
	my.requestUrl = url
	return my
}

// SetMethod 设置请求方法
func (my *Client) SetMethod(method string) *Client {
	my.requestMethod = method
	return my
}

// AddHeaders 设置请求头
func (my *Client) AddHeaders(headers map[string][]string) *Client {
	my.requestHeaders = headers
	return my
}

// SetQueries 设置请求参数
func (my *Client) SetQueries(queries map[string]string) *Client {
	my.requestQueries = queries
	return my
}

// SetAuthorization 设置认证
func (my *Client) SetAuthorization(username, password, title string) *Client {
	my.requestHeaders["Authorization"] = []string{fmt.Sprintf("%s %s", title, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))}

	return my
}

// SetBody 设置请求体
func (my *Client) SetBody(body []byte) *Client {
	my.requestBody = body

	return my
}

// SetJsonBody 设置json请求体
func (my *Client) SetJsonBody(body any) *Client {
	my.SetHeaderContentType("json")

	my.requestBody, my.Err = json.Marshal(body)
	return my
}

// SetXmlBody 设置xml请求体
func (my *Client) SetXmlBody(body any) *Client {
	my.SetHeaderContentType("xml")

	my.requestBody, my.Err = xml.Marshal(body)

	return my
}

// SetFormBody 设置表单请求体
func (my *Client) SetFormBody(body map[string]string) *Client {
	my.SetHeaderContentType("form")

	params := url.Values{}
	for k, v := range body {
		params.Add(k, v)
	}
	my.requestBody = []byte(params.Encode())

	return my
}

// SetFormDataBody 设置表单数据请求体
func (my *Client) SetFormDataBody(texts map[string]string, files map[string]string) *Client {
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
			file, _ := os.Open(v)
			_, e = io.Copy(fileWriter, file)
			if e != nil {
				my.Err = e
				return my
			}
			defer func(file *os.File) {
				e = file.Close()
				if e != nil {
					panic(e)
				}
			}(file)
		}
	}

	my.requestBody = []byte(writer.FormDataContentType())

	return my
}

// SetPlainBody 设置纯文本请求体
func (my *Client) SetPlainBody(text string) *Client {
	my.SetHeaderContentType("plain")

	my.requestBody = []byte(text)

	return my
}

// SetHtmlBody 设置html请求体
func (my *Client) SetHtmlBody(text string) *Client {
	my.SetHeaderContentType("html")

	my.requestBody = []byte(text)

	return my
}

// SetCssBody 设置Css请求体
func (my *Client) SetCssBody(text string) *Client {
	my.SetHeaderContentType("css")

	my.requestBody = []byte(text)

	return my
}

// SetJavascriptBody 设置Javascript请求体
func (my *Client) SetJavascriptBody(text string) *Client {
	my.SetHeaderContentType("javascript")

	my.requestBody = []byte(text)

	return my
}

// SetSteamBody 设置二进制文件
func (my *Client) SetSteamBody(filename string) *Client {
	var (
		err  error
		file *os.File
	)

	file, err = os.Open(filename)
	if err != nil {
		my.Err = err
		return my
	}

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

	my.request.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	my.Err = file.Close()

	return my
}

// SetHeaderContentType 设置请求头内容类型
func (my *Client) SetHeaderContentType(key ContentType) *Client {
	if val, ok := ContentTypes[key]; ok {
		my.requestHeaders["Content-Type"] = []string{val}
	}

	return my
}

// AppendHeaderContentType 追加请求头内容类型
func (my *Client) AppendHeaderContentType(keys ...ContentType) *Client {
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
func (my *Client) SetHeaderAccept(key Accept) *Client {
	if val, ok := Accepts[key]; ok {
		my.requestHeaders["Accept"] = []string{val}
	}

	return my
}

// AppendHeaderAccept 追加请求头接受内容类型
func (my *Client) AppendHeaderAccept(keys ...Accept) *Client {
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
func (my *Client) SetTimeoutSecond(timeoutSecond int64) *Client {
	my.timeoutSecond = timeoutSecond

	return my
}

// GetResponse 获取响应对象
func (my *Client) GetResponse() *http.Response {
	return my.response
}

// ParseByContentType 根据响应头Content-Type自动解析响应体
func (my *Client) ParseByContentType(target any) *Client {
	switch ContentType(my.GetResponse().Header.Get("Content-Type")) {
	case ContentTypeJson:
		my.GetResponseJsonBody(target)
	case ContentTypeXml:
		my.GetResponseXmlBody(target)
	}
	return my
}

// GetResponseRawBody 获取原始响应体
func (my *Client) GetResponseRawBody() []byte {
	return my.responseBody
}

// GetResponseJsonBody 获取json格式响应体
func (my *Client) GetResponseJsonBody(target any) *Client {
	if e := json.Unmarshal(my.responseBody, &target); e != nil {
		my.Err = e
	}
	return my
}

// GetResponseXmlBody 获取xml格式响应体
func (my *Client) GetResponseXmlBody(target any) *Client {
	if e := xml.Unmarshal(my.responseBody, &target); e != nil {
		my.Err = e
	}
	return my
}

// SaveResponseSteamFile 保存二进制到文件
func (my *Client) SaveResponseSteamFile(filename string) *Client {
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
func (my *Client) GetRequest() *http.Request {
	return my.request
}

// GenerateRequest 生成请求对象
func (my *Client) GenerateRequest() *Client {
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

// Send 发送请求
func (my *Client) Send() *Client {
	if !my.isReady {
		my.GenerateRequest()
		if my.Err != nil {
			return my
		}
	}

	my.responseBodyBuffer.Reset() // 重置响应体缓存

	// 发送新的请求
	client := &http.Client{Transport: my.transport}

	// 设置超时
	if my.timeoutSecond > 0 {
		client.Timeout = time.Duration(my.timeoutSecond) * time.Second
	}

	my.response, my.Err = client.Do(my.request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			my.Err = fmt.Errorf("发送失败：%v", err)
		}
	}(my.response.Body)
	if my.Err != nil {
		return my
	}

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
func (my *Client) check() error {
	if my.requestUrl == "" {
		return errors.New("url不能为空")
	}

	if my.requestMethod == "" {
		my.requestMethod = http.MethodGet
	}

	return nil
}

// 设置url参数
func (my *Client) setQueries() {
	if len(my.requestQueries) > 0 {
		queries := url.Values{}
		for k, v := range my.requestQueries {
			queries.Add(k, v)
		}

		my.requestUrl = str.NewBufferByString(my.requestUrl).WriteString("?", queries.Encode()).ToString()
	}
}

// 设置请求头
func (my *Client) addHeaders() {
	for k, v := range my.requestHeaders {
		my.request.Header[k] = append(my.request.Header[k], v...)
	}
}
