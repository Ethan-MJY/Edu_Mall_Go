// Package admin 管理后台API控制器
// 职责: 处理HTTP请求,参数绑定,调用Service层,返回统一响应
package admin

import (
	"mall/adaptor"
	"mall/service/admin"
)

// Ctrl 管理员控制器
type Ctrl struct {
	adaptor adaptor.IAdaptor // 适配器(预留)
	user    *admin.Service    // 管理员业务服务
}

// NewCtrl 创建管理员控制器实例
// 参数: adaptor 适配器,提供数据库和Redis访问
// 返回: Ctrl实例
// 调用链: router.NewRouter -> admin.NewCtrl
func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor), // 初始化业务服务
	}
}
