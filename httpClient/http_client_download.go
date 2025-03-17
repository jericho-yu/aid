package httpClient

import (
	"io"
	"net/http"
	"os"

	"github.com/jericho-yu/aid/array"
	"github.com/jericho-yu/aid/dict"
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
func (my *HttpClientDownload) SaveLocal() *HttpClient {
	defer func() { my.httpClient.isReady = false }()

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
func (my *HttpClientDownload) SendResponse(w http.ResponseWriter, headers ...map[string][]string) *HttpClient {
	defer func() { my.httpClient.isReady = false }()

	client := my.httpClient.beforeSend()
	if my.httpClient.Err != nil {
		return nil
	}

	if my.httpClient.response, my.httpClient.Err = client.Do(my.httpClient.request); my.httpClient.Err != nil {
		return nil
	} else {
		defer my.httpClient.response.Body.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+my.filename)
		w.Header().Set("Content-Type", my.httpClient.response.Header.Get("Content-Type"))

		array.New(headers).Each(func(_ int, header map[string][]string) {
			dict.New(header).Each(func(key string, values []string) {
				array.New(values).Each(func(_ int, value string) { w.Header().Set(key, value) })
			})
		})

		if _, my.httpClient.Err = io.Copy(w, my.httpClient.response.Body); my.httpClient.Err != nil {
			my.httpClient.Err = WriteResponseErr.Wrap(my.httpClient.Err)
		}

		return my.httpClient
	}
}
