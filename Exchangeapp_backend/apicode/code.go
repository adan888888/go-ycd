// Package apicode 统一业务码 code 与 HTTP 状态码的映射。
//
// 约定：
//   - 响应体固定 { code, msg, data }，客户端只看 code
//   - HTTP 状态码由 code 推导，二者一一对应
package apicode

import "net/http"

const (
	CodeOK           = 0    // 成功
	CodeFail         = 1    // 普通业务失败（参数错误、密码错误等）
	CodeUnauthorized = 401  // 未登录 / token 失效
	CodeForbidden    = 403  // 无权限
	CodeNotFound     = 404  // 资源不存在
	CodeYcdExpired   = 2202 // ycd 服务已到期
	CodeServerError  = 500  // 服务器内部错误
)

// Category 业务码分类，供客户端统一处理。
type Category string

const (
	CategorySuccess      Category = "success"
	CategoryBusinessFail Category = "business_fail"
	CategoryUnauthorized Category = "unauthorized"
	CategoryForbidden    Category = "forbidden"
	CategoryYcdExpired   Category = "ycd_expired"
	CategoryServerError  Category = "server_error"
)

// HTTPStatusForCode 由业务码推导 HTTP 状态码。
func HTTPStatusForCode(code int) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeYcdExpired:
		return http.StatusPaymentRequired // 402 订阅/到期
	case CodeServerError:
		return http.StatusInternalServerError
	case CodeNotFound:
		return http.StatusNotFound
	default:
		if code >= 400 && code < 600 {
			return code
		}
		return http.StatusBadRequest
	}
}

// CodeFromHTTPStatus 将旧代码里使用的 HTTP 状态码转为业务码（兼容迁移）。
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
		return CodeFail
	}
}

// CategoryOf 返回业务码所属分类。
func CategoryOf(code int) Category {
	switch code {
	case CodeOK:
		return CategorySuccess
	case CodeUnauthorized:
		return CategoryUnauthorized
	case CodeForbidden:
		return CategoryForbidden
	case CodeYcdExpired:
		return CategoryYcdExpired
	case CodeServerError:
		return CategoryServerError
	default:
		return CategoryBusinessFail
	}
}

// IsGlobalCode 是否需客户端全局拦截（跳转、清 session 等）。
func IsGlobalCode(code int) bool {
	switch code {
	case CodeUnauthorized, CodeForbidden, CodeYcdExpired:
		return true
	default:
		return false
	}
}
