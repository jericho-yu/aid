package httpClient

type (
	ContentTypeStruct struct {
		Key   string
		Value string
	}

	// ContentType http客户端内容类型
	ContentType struct {
		ContentTypeStruct
	}
	Accept struct {
		ContentTypeStruct
	}

	httpClientAccept      = string
	HttpClientContentType = string
)

var (
	ContentTypeJson       HttpClientContentType = "json"
	ContentTypeXml        HttpClientContentType = "xml"
	ContentTypeForm       HttpClientContentType = "form"
	ContentTypeFormData   HttpClientContentType = "form-data"
	ContentTypePlain      HttpClientContentType = "plain"
	ContentTypeHtml       HttpClientContentType = "html"
	ContentTypeCss        HttpClientContentType = "css"
	ContentTypeJavascript HttpClientContentType = "javascript"
	ContentTypeSteam      HttpClientContentType = "steam"

	httpClientContentTypes = []ContentType{
		{ContentTypeStruct{Key: "json", Value: "application/json"}},
		{ContentTypeStruct{Key: "xml", Value: "application/xml"}},
		{ContentTypeStruct{Key: "form", Value: "application/x-www-form-urlencoded"}},
		{ContentTypeStruct{Key: "form-data", Value: "form-data"}},
		{ContentTypeStruct{Key: "plain", Value: "text/plain"}},
		{ContentTypeStruct{Key: "html", Value: "text/html"}},
		{ContentTypeStruct{Key: "css", Value: "text/css"}},
		{ContentTypeStruct{Key: "javascript", Value: "text/javascript"}},
		{ContentTypeStruct{Key: "steam", Value: "application/octet-stream"}},
		{ContentTypeStruct{Key: "any", Value: ""}},
	}

	AcceptJson       httpClientAccept = "json"
	AcceptXml        httpClientAccept = "xml"
	AcceptPlain      httpClientAccept = "plain"
	AcceptHtml       httpClientAccept = "html"
	AcceptCss        httpClientAccept = "css"
	AcceptJavascript httpClientAccept = "javascript"
	AcceptSteam      httpClientAccept = "steam"
	AcceptAny        httpClientAccept = "any"

	httpClientAccepts = []Accept{
		{ContentTypeStruct{Key: "json", Value: "application/json"}},
		{ContentTypeStruct{Key: "xml", Value: "application/xml"}},
		{ContentTypeStruct{Key: "form", Value: ""}},
		{ContentTypeStruct{Key: "plain", Value: "text/plain"}},
		{ContentTypeStruct{Key: "html", Value: "text/html"}},
		{ContentTypeStruct{Key: "css", Value: "text/css"}},
		{ContentTypeStruct{Key: "javascript", Value: "text/javascript"}},
		{ContentTypeStruct{Key: "steam", Value: "application/octet-stream"}},
		{ContentTypeStruct{Key: "any", Value: "*/*"}},
	}
)

// GetValue http客户端请求内容类型
func (ContentType) GetValue(key string) string {
	for _, httpClientContentType := range httpClientContentTypes {
		if httpClientContentType.Key == key {
			return httpClientContentType.Value
		}
	}
	return ""
}

// GetValue http客户端接受内容类型
func (Accept) GetValue(key string) string {
	for _, httpClientAccept := range httpClientAccepts {
		if httpClientAccept.Key == key {
			return httpClientAccept.Value
		}
	}
	return ""
}
