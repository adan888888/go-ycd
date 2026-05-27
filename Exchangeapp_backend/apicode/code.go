// Package apicode 统一业务码 code 定义与辅助判断。
//
// 约定：
//   - 请求到达后端后 HTTP 固定 200
//   - 响应体 { code, msg, data }，code=0 成功，非 0 为业务码
//   - 客户端直接 switch code；需全局拦截见 IsGlobalCode
package apicode

import "net/http"

const (
	CodeOK = 0 // 成功

	// 业务码 1000~1999
	CodeParamInvalid      = 1000 // 参数或业务校验失败
	CodeLoginInvalid      = 1001 // 登录时账号或密码校验失败
	CodeVerifyCodeExpired = 1002 // 验证码已过期
	CodeJsqExpired        = 1003 // jsq（计数器）服务已到期
	CodeNotFound          = 1004 // 资源不存在
	CodeUnauthorized      = 1005 // 未登录 / token 失效
	CodeForbidden         = 1006 // 无权限
	CodeServerError       = 1007 // 服务器内部错误
)

// CodeFromHTTPStatus 内部 scope 等逻辑沿用的 HTTP 语义 → 业务码。
func CodeFromHTTPStatus(httpStatus int) int {
	switch httpStatus {
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusInternalServerError:
		return CodeServerError
	default:
		return CodeParamInvalid
	}
}

func IsSuccess(code int) bool {
	return code == CodeOK
}

// IsGlobalCode 需客户端全局拦截的业务码。
func IsGlobalCode(code int) bool {
	switch code {
	case CodeUnauthorized, CodeForbidden, CodeJsqExpired:
		return true
	default:
		return false
	}
}
