// Package router 路由层
// 职责: HTTP路由注册、中间件配置、服务器生命周期管理
// 路由分组: /api/mall/admin (管理后台) + /api/mall/customer (用户前台)
package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"mall/adaptor"
	"mall/api/admin"
	"mall/api/customer"
	"mall/common"
	"mall/config"
	"net/http"
	"strings"
)

// IRouter 路由器接口
type IRouter interface {
	Register(engine *gin.Engine)              // 注册所有路由
	SpanFilter(r *gin.Context) bool           // 认证过滤器(白名单判断)
	AccessRecordFilter(r *gin.Context) bool   // 访问日志过滤器
}

// Router 路由器结构体
type Router struct {
	FullPPROF bool            // 是否启用pprof性能分析
	rootPath  string          // API根路径: /api/mall
	conf      *config.Config  // 配置对象
	checkFunc func() error    // 健康检查函数(MySQL+Redis连接测试)
	admin     *admin.Ctrl     // 管理后台控制器
	customer  *customer.Ctrl  // 用户前台控制器
}

// NewRouter 创建路由器实例
// 参数:
//   - conf: 配置对象
//   - adaptor: 适配器(提供数据库、Redis访问)
//   - checkFunc: 健康检查函数
// 返回: Router实例
// 调用链: main.main -> NewRouter
func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),      // 初始化管理后台控制器
		customer:  customer.NewCtrl(adaptor),   // 初始化用户前台控制器
	}
}

// checkServer 健康检查接口处理函数
// 路由: GET/POST /ping
// 返回: MySQL和Redis连接状态
func (r *Router) checkServer() func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// Register 注册所有路由
// 执行流程:
//   1. 注册pprof路由(如果启用)
//   2. 注册/ping健康检查
//   3. 注册业务路由(admin + customer)
// 调用链: main.main -> NewApp -> Register
func (r *Router) Register(app *gin.Engine) {
	// 注册pprof性能分析工具
	if r.conf.Server.EnablePprof {
		SetupPprof(app, "/debug/pprof")
	}

	// 健康检查接口
	app.Any("/ping", r.checkServer())

	// 业务路由组: /api/mall
	root := app.Group(r.rootPath)
	r.route(root)
}

// SpanFilter 认证白名单过滤器
// 参数: ctx Gin上下文
// 返回: true表示需要认证,false表示跳过认证(白名单)
// 白名单: 登录、验证码等接口无需认证
func (r *Router) SpanFilter(ctx *gin.Context) bool {
	path := strings.Replace(ctx.Request.URL.Path, r.rootPath, "", 1)
	_, ok := AdminAuthWhiteList[path]
	if ok {
		return false // 在白名单中,跳过认证
	}
	return true // 不在白名单,需要认证
}

// AccessRecordFilter 访问日志过滤器
// 参数: ctx Gin上下文
// 返回: true表示记录日志,false表示跳过
// 当前实现: 所有请求均记录日志
func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

// route 注册业务路由
// 分组:
//   - /api/mall/customer: 用户前台API
//   - /api/mall/admin: 管理后台API
func (r *Router) route(root *gin.RouterGroup) {
	r.customerRoute(root)
	r.adminRoute(root)
}

// customerRoute 注册用户前台路由
// 路由前缀: /api/mall/customer
// 认证: AuthMiddleware(用户Token)
// 白名单: 通过SpanFilter判断
// TODO: 完善JWT Token解析逻辑
func (r *Router) customerRoute(root *gin.RouterGroup) {
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.User, error) {
		// TODO: 实现真实的JWT Token解析
		return &common.User{}, nil
	}))
	// 用户信息接口
	cstRoot.Any("/user/info", r.admin.GetUserInfo)
}

// adminRoute 注册管理后台路由
// 路由前缀: /api/mall/admin
// 认证: AdminAuthMiddleware(管理员Token)
// 白名单: 登录、验证码等接口无需认证
// TODO: 完善JWT Token解析逻辑
func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.AdminUser, error) {
		// TODO: 实现真实的JWT Token解析
		return &common.AdminUser{
			UserID: 1,
			Name:   "admin",
		}, nil
	}))

	// ========== 登录相关(无需认证,在白名单中) ==========
	// 获取滑块验证码
	adminRoot.GET("/v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	// 校验滑块验证码
	adminRoot.POST("/v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)

	// ========== 用户管理(需要认证) ==========
	// 获取用户信息
	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	// 创建用户
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	// 更新用户
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
}
