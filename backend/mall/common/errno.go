// Package common 公共组件模块
// 职责: 统一错误码和响应码定义,提供错误处理方法
package common

// Errno 统一错误码结构体
// Code: 错误码,200表示成功
// Msg: 错误消息,返回给客户端
// ErrMsg: 详细错误信息,用于日志记录
type Errno struct {
	Code   int
	Msg    string
	ErrMsg string
}

// Error 实现error接口
func (err Errno) Error() string {
	return err.Msg
}

// WithMsg 追加错误消息
// 用于在原始错误信息基础上添加上下文信息
func (err Errno) WithMsg(msg string) Errno {
	err.Msg = err.Msg + "," + msg
	return err
}

// WithErr 追加原始错误信息
// 将底层错误信息附加到ErrMsg,用于日志记录
func (err Errno) WithErr(rawErr error) Errno {
	var msg string
	if rawErr != nil {
		msg = rawErr.Error()
	}
	err.ErrMsg = err.Msg + "," + msg
	return err
}

// IsOk 判断是否成功
// 返回Code是否为200
func (err Errno) IsOk() bool {
	return err.Code == 200
}

// 预定义错误码
var (
	// HTTP标准错误码
	OK            = Errno{Code: 200, Msg: "OK"}
	ServerErr     = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Param Error"}
	AuthErr       = Errno{Code: 401, Msg: "Auth Error"}
	PermissionErr = Errno{Code: 403, Msg: "Permission Error"}

	// 基础设施错误码 (10000-10999)
	DatabaseErr = Errno{Code: 10000, Msg: "Database Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	// 业务错误码 (11000+)
	UserNotFoundErr   = Errno{Code: 11001, Msg: "User Not Found"}
	InvalidCaptchaErr = Errno{Code: 11002, Msg: "滑块校验失败，请重试"}
)
