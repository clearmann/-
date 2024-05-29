package g

import "fmt"

const (
	SUCCESS = 0
	FAIL    = 500
)

// Result 自定义业务 error
type Result struct {
	code    int
	message string
}

func (r Result) Code() int {
	return r.code
}
func (r Result) Message() string {
	return r.message
}

var (
	_code    = make(map[int]bool)   // 注册过的业务状态码，防止重复
	_message = make(map[int]string) // 根据业务状态码，得到 错误信息
)

// 注册一个业务码，不允许重复注册
func RegisterResult(code int, message string) Result {
	if _code[code] {
		panic(fmt.Sprintf("业务状态码 %d 已经存在，请更换一个", code))
	}
	if len(message) == 0 {
		panic("错误信息不能为空")
	}
	return Result{
		code:    code,
		message: message,
	}
}

// 根据业务状态码 获取 错误信息
func GetMsg(code int) string {
	return _message[code]
}

var (
	OkResult   = RegisterResult(SUCCESS, "OK")
	FailResult = RegisterResult(FAIL, "FAIL")
)
