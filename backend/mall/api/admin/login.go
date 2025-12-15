// Package admin 管理后台API控制器-登录模块
// 职责: 验证码相关接口处理
package admin

import (
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/common"
	"mall/service/dto"
)

// GetSmsCodeCaptcha 获取滑块验证码接口
// 路由: GET /api/mall/admin/v1/user/verify/captcha
// 参数: 无
// 返回: 验证码图片Base64、Key、滑块尺寸等信息
// 白名单: 无需Token认证
// 调用链: router -> GetSmsCodeCaptcha -> service.GetSlideCaptcha
func (c *Ctrl) GetSmsCodeCaptcha(ctx *gin.Context) {
	// 1. 参数绑定(Query参数)
	req := &dto.GetVerifyCaptchaReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	// 2. 调用Service层获取验证码
	resp, errno := c.user.GetSlideCaptcha(ctx.Request.Context())

	// 3. 返回响应
	api.WriteResp(ctx, resp, errno)
}

// CheckSmsCodeCaptcha 校验滑块验证码接口
// 路由: POST /api/mall/admin/v1/user/verify/captcha/check
// 参数: JSON Body - Key(验证码标识) + SlideX/SlideY(用户滑动坐标)
// 返回: Ticket(验证通过凭证,有效期5分钟)
// 白名单: 无需Token认证
// 用途: 验证通过后返回Ticket,用于后续登录接口
// 调用链: router -> CheckSmsCodeCaptcha -> service.CheckSlideCaptcha
func (c *Ctrl) CheckSmsCodeCaptcha(ctx *gin.Context) {
	// 1. 参数绑定(JSON Body)
	req := &dto.CheckCaptchaReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}

	// 2. 调用Service层校验验证码
	resp, errno := c.user.CheckSlideCaptcha(ctx.Request.Context(), req)

	// 3. 返回响应
	api.WriteResp(ctx, resp, errno)
}
