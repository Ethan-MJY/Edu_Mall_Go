// Package consts 常量定义模块
// 职责: 定义系统级常量,包括Token Key、状态常量等
package consts

const (
	AdminTokenKey   = "token"          // 管理员Token在请求头中的键名
	UserTokenKey    = "token"          // 用户Token在请求头中的键名
	CustomerUserKey = "user_key"       // 客户端用户信息在Context中的键名
	AdminUserKey    = "admin_user_key" // 管理员用户信息在Context中的键名
)

const (
	IsEnable  = 1  // 启用状态
	IsDisable = -1 // 禁用状态
)
