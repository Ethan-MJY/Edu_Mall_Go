// Package router 路由层-性能分析工具
// 职责: 集成pprof性能分析工具
package router

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

// SetupPprof 设置pprof性能分析路由
// 参数:
//   - r: Gin引擎
//   - prefix: 路由前缀,通常为"/debug/pprof"
// 功能: 注册pprof各个性能分析端点
// 使用: 通过浏览器访问 http://localhost:port/debug/pprof 查看性能数据
// 端点说明:
//   - /: 性能分析首页
//   - /cmdline: 命令行参数
//   - /profile: CPU性能分析(30秒采样)
//   - /symbol: 符号解析
//   - /trace: 请求追踪
//   - /heap: 堆内存分析
//   - /goroutine: 协程栈信息
//   - /block: 阻塞分析
//   - /mutex: 互斥锁分析
func SetupPprof(r *gin.Engine, prefix string) {
	group := r.Group(prefix)
	{
		group.GET("/", gin.WrapF(pprof.Index))              // 性能分析首页
		group.GET("/cmdline", gin.WrapF(pprof.Cmdline))     // 命令行参数
		group.GET("/profile", gin.WrapF(pprof.Profile))     // CPU性能分析
		group.GET("/symbol", gin.WrapF(pprof.Symbol))       // 符号解析
		group.GET("/trace", gin.WrapF(pprof.Trace))         // 请求追踪
		group.GET("/heap", pprofHandler("heap"))            // 堆内存分析
		group.GET("/goroutine", pprofHandler("goroutine"))  // 协程栈信息
		group.GET("/block", pprofHandler("block"))          // 阻塞分析
		group.GET("/mutex", pprofHandler("mutex"))          // 互斥锁分析
	}
}

// pprofHandler 统一返回pprof的Profile数据
// 参数: name 性能分析类型(heap/goroutine/block/mutex)
// 返回: Gin处理函数
func pprofHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		pprof.Handler(name).ServeHTTP(c.Writer, c.Request)
	}
}
