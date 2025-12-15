// Package router 路由层-认证中间件
// 职责: JWT Token认证中间件实现
package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"mall/common"
	"mall/consts"
	"net/http"
)

// TokenFun 用户Token解析函数类型
// 参数: context和token字符串
// 返回: 用户对象和错误
type TokenFun func(ctx context.Context, token string) (*common.User, error)

// TokenAdminFun 管理员Token解析函数类型
// 参数: context和token字符串
// 返回: 管理员用户对象和错误
type TokenAdminFun func(ctx context.Context, token string) (*common.AdminUser, error)

// AuthMiddleware 用户侧认证中间件
// 参数:
//   - filter: 白名单过滤器,返回false则跳过认证
//   - getTokenFun: Token解析函数
// 返回: Gin中间件函数
// 功能:
//   1. 检查是否在白名单中
//   2. 从Header中获取Token
//   3. 解析Token获取用户信息
//   4. 将用户信息存入Context
// Header: user_key
func AuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 白名单检查,如果在白名单中,直接跳过认证
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		// 从Header中获取Token
		token := ctx.GetHeader(consts.UserTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		// 解析Token获取用户信息
		user, err := getTokenFun(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr.WithErr(err))
			ctx.Abort()
			return
		}

		// 将用户信息存入Context,供后续Handler使用
		ctx.Set(consts.CustomerUserKey, user)
		ctx.Next()
	}
}

// AdminAuthMiddleware 管理后台认证中间件
// 参数:
//   - filter: 白名单过滤器,返回false则跳过认证
//   - getTokenFun: Token解析函数
// 返回: Gin中间件函数
// 功能:
//   1. 检查是否在白名单中(登录、验证码等接口)
//   2. 从Header中获取Token
//   3. 解析Token获取管理员信息
//   4. 将管理员信息存入Context
// Header: token
func AdminAuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenAdminFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 白名单检查,如果在白名单中,直接跳过认证
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		// 从Header中获取Token
		token := ctx.GetHeader(consts.AdminTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}

		// 解析Token获取管理员信息
		user, err := getTokenFun(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr.WithErr(err))
			ctx.Abort()
			return
		}

		// 将管理员信息存入Context,供后续Handler使用
		ctx.Set(consts.AdminUserKey, user)
		ctx.Next()
	}
}
