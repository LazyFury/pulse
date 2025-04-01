package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazyfury/pulse/framework/constants"
	"github.com/lazyfury/pulse/framework/global"
)

type ApiJsonResponse struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Data     interface{}    `json:"data"`
	HttpCode int            `json:"-"`
	Extra    map[string]any `json:"extra"`
}

type HtmlTemplateResponse struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Data     interface{}    `json:"data"`
	HttpCode int            `json:"-"`
	Extra    map[string]any `json:"extra"`
	Template string         `json:"template"`
}

func Html(code int, message string, data interface{}, template string, httpCode int) *HtmlTemplateResponse {
	return &HtmlTemplateResponse{
		Code:     code,
		Message:  message,
		Data:     data,
		HttpCode: httpCode,
		Template: template,
	}
}

// withExtra
func (a *HtmlTemplateResponse) WithExtra(extra map[string]any) *HtmlTemplateResponse {
	if a.Extra == nil {
		a.Extra = make(map[string]any)
	}
	for k, v := range extra {
		a.Extra[k] = v
	}
	return a
}

func (a *HtmlTemplateResponse) Done(ctx *gin.Context) {
	a.WithExtra(map[string]any{
		"path":    ctx.Request.URL.Path,
		"request": ctx.Request,
	})

	namespace := global.GetConfig().GetString(constants.THEME)
	if namespace == "" {
		namespace = "default"
	}
	ctx.HTML(a.HttpCode, fmt.Sprintf("%s/%s", namespace, a.Template), a)
}

func NewApiJsonResponse(code int, message string, data interface{}, httpCode int) *ApiJsonResponse {
	return &ApiJsonResponse{
		Code:     code,
		Message:  message,
		Data:     data,
		HttpCode: httpCode,
	}
}

func (a *ApiJsonResponse) ToJson() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// with ctx
func (a *ApiJsonResponse) Done(ctx *gin.Context) {
	ctx.JSON(a.HttpCode, a)
}

func (a *ApiJsonResponse) WithExtra(extra map[string]any) *ApiJsonResponse {
	if a.Extra == nil {
		a.Extra = make(map[string]any)
	}
	for k, v := range extra {
		a.Extra[k] = v
	}
	return a
}

// with data
func (a *ApiJsonResponse) WithData(data interface{}) *ApiJsonResponse {
	a.Data = data
	return a
}

// success
func Success(data interface{}) *ApiJsonResponse {
	return NewApiJsonResponse(http.StatusOK, "ok", data, http.StatusOK)
}

// error
func Error(error ApiError, message string) *ApiJsonResponse {
	if message == "" {
		message = error.GetMsgFromErrCode()
	}
	var httpCode int = http.StatusBadRequest
	if error.IsHttpCode() {
		httpCode = error.Code
	}
	return NewApiJsonResponse(error.Code, message, nil, httpCode)
}

// not found
func NotFound(message string) *ApiJsonResponse {
	return NewApiJsonResponse(http.StatusNotFound, message, nil, http.StatusNotFound)
}

// bad request
func BadRequest(message string) *ApiJsonResponse {
	return NewApiJsonResponse(http.StatusBadRequest, message, nil, http.StatusBadRequest)
}
