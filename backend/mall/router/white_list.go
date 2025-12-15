// Package router 路由层-白名单配置
// 职责: 定义无需认证的接口白名单
package router

// AdminAuthWhiteList 管理后台认证白名单
// key: 接口路径(不包含/api/mall前缀)
// value: true表示在白名单中,无需Token认证
// 白名单接口:
//   - /ping: 健康检查
//   - /metrics: 监控指标
//   - /admin/v1/user/verify/*: 验证码相关接口
//   - /admin/v1/user/mobile/*: 手机号登录接口
//   - /admin/v1/user/password/reset: 密码重置
var AdminAuthWhiteList = map[string]bool{
	"/ping":                                true, // 健康检查
	"/metrics":                             true, // 监控指标
	"/admin/v1/user/verify/captcha/check":  true, // 滑块验证码校验
	"/admin/v1/user/verify/captcha":        true, // 获取滑块验证码
	"/admin/v1/user/verify/smscode":        true, // 获取短信验证码
	"/admin/v1/user/mobile/verify_login":   true, // 手机号验证码登录
	"/admin/v1/user/mobile/password_login": true, // 手机号密码登录
	"/admin/v1/user/password/reset":        true, // 密码重置
}
