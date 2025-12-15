// Package router 路由层-访问日志中间件
// 职责: 记录HTTP请求的详细访问日志
package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mall/consts"
	"mall/utils/logger"
	"time"
)

// GetRequestBody 获取请求Body内容
// 参数: ctx Gin上下文
// 返回: Body字符串
func GetRequestBody(ctx *gin.Context) string {
	data, _ := io.ReadAll(ctx.Request.Body)
	return string(data)
}

// GetResponseBody 获取响应Body内容
// 参数: ctx Gin上下文
// 返回: Body字符串
// 注意: 需要配合responseWriterWrapper使用
func GetResponseBody(ctx *gin.Context) string {
	resp := ctx.Request.Response
	if resp == nil || resp.Body == nil {
		return ""
	}
	data, _ := io.ReadAll(ctx.Request.Response.Body)
	return string(data)
}

// responseWriterWrapper 响应Writer包装器
// 用途: 拦截响应内容,用于日志记录
type responseWriterWrapper struct {
	gin.ResponseWriter        // 嵌入原始ResponseWriter
	Writer             io.Writer // 多路写入器,同时写入原始Writer和Buffer
}

// Write 实现io.Writer接口
// 将响应内容同时写入原始Writer和Buffer
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// AccessLogMiddleware 访问日志中间件
// 参数: filter 过滤器函数,返回false则跳过日志记录
// 返回: Gin中间件函数
// 功能:
//   1. 记录请求信息(IP、方法、路径、参数、Body、Token)
//   2. 拦截响应内容
//   3. 记录响应信息(状态码、响应Body、耗时)
//   4. 输出到日志系统
// 日志字段:
//   - ip: 客户端IP
//   - method: HTTP方法
//   - path: 请求路径
//   - params: Query参数
//   - body: 请求Body
//   - token: 用户Token
//   - status: 响应状态码
//   - resp: 响应Body(最多1024字符)
//   - dur_ms: 耗时(毫秒)
func AccessLogMiddleware(filter func(*gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 过滤器判断,如果返回false则跳过日志记录
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		// 读取请求Body并重新设置(因为Body只能读一次)
		body := GetRequestBody(ctx)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))

		// 记录请求开始时间
		begin := time.Now()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("ip", ctx.RemoteIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("params", ctx.Request.URL.RawQuery),
			zap.Any("body", body),
			zap.String("token", ctx.GetHeader(consts.UserTokenKey)),
		}

		// 包装ResponseWriter,拦截响应内容
		var responseBody bytes.Buffer
		multiWriter := io.MultiWriter(ctx.Writer, &responseBody)
		ctx.Writer = &responseWriterWrapper{
			ResponseWriter: ctx.Writer,
			Writer:         multiWriter,
		}

		// 执行后续Handler
		ctx.Next()

		// 获取响应Body(限制最大1024字符)
		respBody := responseBody.String()
		if len(respBody) > 1024 {
			respBody = respBody[:1024]
		}

		// 追加响应信息到日志字段
		fields = append(fields, zap.Int64("dur_ms", time.Since(begin).Milliseconds()))
		fields = append(fields, zap.Int("status", ctx.Writer.Status()))
		fields = append(fields, zap.String("resp", respBody))

		// 输出访问日志
		logger.Info("access_log", fields...)
	}
}
