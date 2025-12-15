// Package api API控制器基础层
// 职责: 统一响应格式、用户上下文获取
package api

import (
	"github.com/gin-gonic/gin"
	"mall/common"
	"mall/consts"
	"net/http"
)

// Resp 统一响应结构体
type Resp struct {
	Code   int    `json:"code"`    // 错误码,200表示成功
	Msg    string `json:"msg"`     // 错误消息,返回给客户端
	ErrMsg string `json:"err_msg"` // 详细错误信息,用于调试
	Data   any    `json:"data"`    // 响应数据
}

// WriteResp 写入统一响应
// 参数:
//   - ctx: Gin上下文
//   - data: 响应数据
//   - errno: 错误码对象
// 功能: 将数据和错误码封装为统一格式并返回
func WriteResp(ctx *gin.Context, data any, errno common.Errno) {
	ctx.JSON(http.StatusOK, Resp{
		Code:   errno.Code,
		Msg:    errno.Msg,
		ErrMsg: errno.ErrMsg,
		Data:   data,
	})
}

// GetUserFromCtx 从Context获取用户信息
// 参数: ctx Gin上下文
// 返回: 用户对象指针,不存在返回nil
// 用途: 在需要用户信息的Handler中调用
func GetUserFromCtx(ctx *gin.Context) *common.User {
	user, exist := ctx.Get(consts.CustomerUserKey)
	if !exist {
		return nil
	}
	return user.(*common.User)
}

// GetAdminUserFromCtx 从Context获取管理员信息
// 参数: ctx Gin上下文
// 返回: 管理员对象指针,不存在返回nil
// 用途: 在需要管理员信息的Handler中调用(记录操作人等)
func GetAdminUserFromCtx(ctx *gin.Context) *common.AdminUser {
	user, exist := ctx.Get(consts.AdminUserKey)
	if !exist {
		return nil
	}
	return user.(*common.AdminUser)
}
