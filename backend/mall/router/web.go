// Package router 路由层-HTTP服务器
// 职责: HTTP服务器启动和优雅关闭
package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"mall/utils/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// App HTTP服务器应用
type App struct {
	server *gin.Engine // Gin引擎
	addr   string       // 监听地址,格式":port"
}

// NewApp 创建HTTP服务器应用实例
// 参数:
//   - port: 监听端口
//   - router: 路由器,负责注册所有路由
// 返回: App实例
// 配置:
//   - Gin运行模式: ReleaseMode
//   - 中间件: Recovery(全局panic恢复) + AccessLog(访问日志)
// 调用链: main.main -> NewApp
func NewApp(port int, router IRouter) *App {
	// 设置为生产模式,减少日志输出
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Recover中间件,全局捕获panic,防止程序崩溃
	engine.Use(gin.Recovery())

	// 访问日志中间件,记录每个请求的详细信息
	// 支持自定义过滤器,某些接口可以不记录日志
	engine.Use(AccessLogMiddleware(router.AccessRecordFilter))

	// 注册业务路由
	router.Register(engine)

	return &App{
		server: engine,
		addr:   ":" + strconv.Itoa(port),
	}
}

// Run 启动HTTP服务器并等待优雅关闭
// 功能:
//   1. 异步启动HTTP服务器
//   2. 监听系统信号(SIGINT/SIGTERM)
//   3. 收到信号后优雅关闭服务器(等待5秒)
// 调用链: main.main -> app.Run
func (app *App) Run() {
	srv := &http.Server{
		Addr:    app.addr,
		Handler: app.server,
	}

	// 异步启动HTTP服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %v", err)
		}
	}()

	logger.Debug(fmt.Sprintf("server started, endpoint: http://localhost%s", app.addr))

	// 监听系统信号,实现优雅关闭
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	msg := <-closeCh // 阻塞等待信号

	logger.Warn("server closing: ", zap.String("msg", msg.String()))

	// 优雅关闭:等待5秒让现有请求处理完成
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
