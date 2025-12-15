// Package admin 管理员业务逻辑层
// 职责: 实现管理员相关的业务逻辑
// 依赖: adminUser(数据访问) + verify(验证码Redis) + captcha(滑块验证码)
package admin

import (
	"github.com/wenlng/go-captcha/v2/slide"
	"mall/adaptor"
	"mall/adaptor/redis"
	"mall/adaptor/repo/admin"
	"mall/utils/captcha"
)

// Service 管理员服务结构体
type Service struct {
	adminUser admin.IAdminUser // 管理员用户数据访问接口
	verify    redis.IVerify    // 验证码Redis操作接口
	captcha   slide.Captcha    // 滑块验证码生成器
}

// NewService 创建管理员服务实例
// 参数: adaptor 适配器,提供数据库和Redis访问
// 返回: Service实例
// 调用链: api.NewCtrl -> NewService
func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUser(adaptor),   // 初始化用户数据访问
		verify:    redis.NewVerify(adaptor),      // 初始化验证码Redis操作
		captcha:   captcha.NewSlideCaptcha(),     // 初始化滑块验证码生成器
	}
}
