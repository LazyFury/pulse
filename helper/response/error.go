package response

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrCodeNeedBindingPhone int = iota + 1000 // 需要绑定手机号
	ErrMpWxBinding                            // 微信绑定
	ErrCodeForbidden                          // 已封禁
	ErrTipsFeedback                           // 提交反馈
)

var errCodeMaps = map[int]string{
	ErrCodeNeedBindingPhone: "需要绑定手机号",
	ErrMpWxBinding:          "微信绑定",
	ErrCodeForbidden:        "已封禁",
	ErrTipsFeedback:         "提交反馈",
}

var httpCodeArray = []int{200, 400, 401, 403, 404, 500}

// 通用
var ErrNotFound = ApiError{404, "Not Found"}
var ErrInternalError = ApiError{500, "Internal Error"}
var ErrInvalidArgument = ApiError{400, "Invalid Argument"}
var ErrNotAuthorized = ApiError{401, "Not Authorized"}
var ErrBadRequest = ApiError{400, "Bad Request"}
var ErrForbidden = ApiError{403, "Forbidden"}

// ok
var ErrOk = ApiError{200, "Ok"}

// 自定义 code
var ErrNotAuth = ApiError{401, "Not Auth"}
var ErrNotPermission = ApiError{403, "Not Permission"}
var ErrNeedBindingPhone = ApiError{ErrCodeNeedBindingPhone, "Binding Phone"}

// 创建新的错误码
func NewError(code int, message string) ApiError {
	return ApiError{code, message}
}

// Error returns the error message of ApiError.
func (e ApiError) Error() string {
	return e.Message
}

// is http code
func (e ApiError) IsHttpCode() bool {
	for _, v := range httpCodeArray {
		if v == e.Code {
			return true
		}
	}
	return false
}

// GetMsgFromErrCode returns the message of the error code.
// If the error code is not in errCodeMaps, it returns the original message.
func (e ApiError) GetMsgFromErrCode() string {
	var str string
	if v, ok := errCodeMaps[e.Code]; ok {
		str = v
	} else {
		str = e.Error()
	}
	return str
}
