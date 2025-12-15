// Package admin 管理后台API控制器-用户管理
// 职责: 管理员用户CRUD接口处理
package admin

import (
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/common"
	"mall/service/dto"
)

// GetUserInfo 获取管理员用户信息接口
// 路由: GET /api/mall/admin/v1/user/info
// 参数: 无(从Token解析当前用户)
// 返回: 用户ID、姓名等基本信息
// 认证: 需要Token
// 调用链: router -> GetUserInfo -> service.GetUserInfo -> repo.GetUserInfo
func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	// 1. 从Context获取当前登录用户
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}

	// 2. 调用Service层获取用户信息
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})

	// 3. 返回响应
	api.WriteResp(ctx, resp, errno)
}

// CreateUser 创建管理员用户接口
// 路由: POST /api/mall/admin/v1/user/create
// 参数: JSON Body - Name(姓名)、NickName(昵称)、Mobile(手机号)、Sex(性别)
// 返回: 新用户ID
// 认证: 需要Token
// 权限: 需要用户管理权限(TODO)
// 调用链: router -> CreateUser -> service.CreateUser -> repo.CreateUser
func (c *Ctrl) CreateUser(ctx *gin.Context) {
	// 1. 从Context获取当前登录用户
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}

	// 2. 参数绑定(JSON Body)
	req := &dto.CreateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}

	// 3. 调用Service层创建用户
	userId, errno := c.user.CreateUser(ctx.Request.Context(), user, req)

	// 4. 返回响应(包含新用户ID)
	api.WriteResp(ctx, map[string]int64{
		"id": userId,
	}, errno)
}

// UpdateUser 更新管理员用户信息接口
// 路由: POST /api/mall/admin/v1/user/update
// 参数: JSON Body - ID(用户ID)、Name(姓名)、NickName(昵称)、Sex(性别)
// 返回: 无
// 认证: 需要Token
// 权限: 需要用户管理权限(TODO)
// 可更新字段: 姓名、昵称、性别
// 调用链: router -> UpdateUser -> service.UpdateUser -> repo.UpdateUser
func (c *Ctrl) UpdateUser(ctx *gin.Context) {
	// 1. 从Context获取当前登录用户
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}

	// 2. 参数绑定(JSON Body)
	req := &dto.UpdateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}

	// 3. 调用Service层更新用户
	errno := c.user.UpdateUser(ctx.Request.Context(), user, req)

	// 4. 返回响应
	api.WriteResp(ctx, nil, errno)
}

// UpdateUserStatus 更新管理员用户状态接口
// 路由: POST /api/mall/admin/v1/user/status
// 参数: JSON Body - ID(用户ID)、Status(状态: 1启用/-1禁用)
// 返回: 无
// 认证: 需要Token
// 权限: 需要用户管理权限(TODO)
// 用途: 启用或停用管理员账号
// 调用链: router -> UpdateUserStatus -> service.UpdateUserStatus -> repo.UpdateUserStatus
func (c *Ctrl) UpdateUserStatus(ctx *gin.Context) {
	// 1. 从Context获取当前登录用户
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}

	// 2. 参数绑定(JSON Body)
	req := &dto.UpdateUserStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
		return
	}

	// 3. 调用Service层更新状态
	errno := c.user.UpdateUserStatus(ctx.Request.Context(), user, req)

	// 4. 返回响应
	api.WriteResp(ctx, nil, errno)
}
