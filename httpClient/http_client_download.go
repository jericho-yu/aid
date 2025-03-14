package httpClient

import (
	"io"
	"os"

	processBar "github.com/schollz/progressbar/v3"
)

type HttpClientDownload struct {
	httpClient     *HttpClient
	filename       string
	processContent string
}

var HttpClientDownloadApp HttpClientDownload

// New 实例化http客户端下载器
func (*HttpClientDownload) New(httpClient *HttpClient, filename string) *HttpClientDownload {
	return &HttpClientDownload{httpClient: httpClient, filename: filename}
}

// SetProcessContent 设置终端进度条标题
func (my *HttpClientDownload) SetProcessContent(processContent string) *HttpClientDownload {
	my.processContent = processContent

	return my
}

// Save 保存到本地
func (my *HttpClientDownload) Save() *HttpClient {
	client := my.httpClient.beforeSend()
	if my.httpClient.Err != nil {
		return my.httpClient
	}

	if my.httpClient.response, my.httpClient.Err = client.Do(my.httpClient.request); my.httpClient.Err != nil {
		return my.httpClient
	} else {
		defer my.httpClient.response.Body.Close()

		f, _ := os.OpenFile(my.filename, os.O_RDWR|os.O_CREATE, 0644)
		defer f.Close()

		if my.processContent != "" {
			_, _ = io.Copy(io.MultiWriter(f, processBar.DefaultBytes(my.httpClient.response.ContentLength, my.processContent)), my.httpClient.response.Body)
		} else {
			_, _ = io.Copy(f, my.httpClient.response.Body)
		}

		return my.httpClient
	}
}

// Send 发送到客户端
func (my *HttpClientDownload) Send() *HttpClient {
	return my.httpClient
}
