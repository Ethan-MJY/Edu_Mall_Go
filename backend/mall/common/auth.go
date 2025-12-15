// Package common 公共组件模块
// 职责: 定义用户认证相关数据结构
package common

// AdminUser 管理员用户信息
// 用于认证中间件解析Token后存储到Context
type AdminUser struct {
	UserID int64  `json:"user_id"` // 管理员ID
	Name   string `json:"name"`    // 管理员姓名
}

// User 前台用户信息
// 用于认证中间件解析Token后存储到Context
type User struct {
	UserID   int64  `json:"user_id"`   // 用户ID
	NickName string `json:"nick_name"` // 用户昵称
}
