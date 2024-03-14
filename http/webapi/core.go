package webapi

var RespDefault func() *WebHTTPApi = func() *WebHTTPApi {
	return NewWebHTTPApi(200)
}

type IWebHTTPApi interface {
	Code() int
	Message() string
	Data() interface{}
	Extra() map[string]interface{}
	Page() int64
	Limit() int64
	Total() int64
	HttpCode() int
}

type WebHTTPApi struct {
	code     int
	message  string
	data     interface{}
	extra    map[string]interface{}
	page     int64
	limit    int64
	total    int64
	httpCode int
}

// NewWebHTTPApi create a new WebHTTPApi for user quickly create a response
// code: response "code" in body
func NewWebHTTPApi(code int) *WebHTTPApi {
	return &WebHTTPApi{
		code: code,
	}
}

func (w *WebHTTPApi) Code() int                     { return w.code }
func (w *WebHTTPApi) Message() string               { return w.message }
func (w *WebHTTPApi) Data() interface{}             { return w.data }
func (w *WebHTTPApi) Extra() map[string]interface{} { return w.extra }
func (w *WebHTTPApi) Page() int64                   { return w.page }
func (w *WebHTTPApi) Limit() int64                  { return w.limit }
func (w *WebHTTPApi) Total() int64                  { return w.total }
func (w *WebHTTPApi) HttpCode() int                 { return w.httpCode }

func (w *WebHTTPApi) SetCode(code int) *WebHTTPApi {
	w.code = code
	return w
}

func (w *WebHTTPApi) SetMessage(message string) *WebHTTPApi {
	w.message = message
	return w
}

func (w *WebHTTPApi) SetData(data interface{}) *WebHTTPApi {
	w.data = data
	return w
}

func (w *WebHTTPApi) SetExtra(extra map[string]interface{}) *WebHTTPApi {
	w.extra = extra
	return w
}

func (w *WebHTTPApi) SetPage(page int64) *WebHTTPApi {
	w.page = page
	return w
}

func (w *WebHTTPApi) SetLimit(limit int64) *WebHTTPApi {
	w.limit = limit
	return w
}

func (w *WebHTTPApi) SetTotal(total int64) *WebHTTPApi {
	w.total = total
	return w
}

func (w *WebHTTPApi) SetHttpCode(httpCode int) *WebHTTPApi {
	w.httpCode = httpCode
	return w
}

func WebHttpApiResponseHandler(resp interface{}) (code int, data map[string]interface{}) {
	v, ok := resp.(*WebHTTPApi)
	var defaultHttpCord = 200
	if ok {
		if v.HttpCode() != 0 {
			defaultHttpCord = v.HttpCode()
		}
		v.HttpCode()

		return 200, map[string]interface{}{
			"code":    v.Code(),
			"message": v.Message(),
			"data":    v.Data(),
			"extra":   v.Extra(),
			"page":    v.Page(),
			"limit":   v.Limit(),
			"total":   v.Total(),
		}
	}

	return defaultHttpCord, map[string]interface{}{
		"code": 200,
		"data": resp,
	}
}
