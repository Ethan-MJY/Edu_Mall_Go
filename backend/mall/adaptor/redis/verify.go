// Package redis Redis操作层-验证码模块
// 职责: 封装验证码相关的Redis存储操作
// 存储内容: 滑块验证码的Key和Ticket
package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"mall/adaptor"
	"mall/config"
	"time"
)

// IVerify 验证码Redis操作接口
// 提供验证码Key和Ticket的存取操作
type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error    // 存储验证码Key
	GetCaptchaKey(ctx context.Context, key string) (string, error)                              // 获取验证码Key(获取后删除)
	SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error // 存储验证码Ticket
	GetCaptchaTicket(ctx context.Context, key string) (string, error)                           // 获取验证码Ticket(获取后删除)
}

// Verify 验证码Redis操作实现
type Verify struct {
	redis *redis.Client // Redis客户端
}

// NewVerify 创建验证码Redis操作实例
// 参数: adaptor 适配器,提供Redis连接
// 返回: Verify实例
// 调用链: service.NewService -> NewVerify
func NewVerify(adaptor adaptor.IAdaptor) *Verify {
	return &Verify{
		redis: adaptor.GetRedis(),
	}
}

// fmtVerifyCaptchaKey 格式化验证码Key的Redis键名
// 格式: <服务名>:captcha:<key>
// 示例: edu_mall:captcha:abc123
func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerFullName, key)
}

// fmtVerifyCaptchaTicket 格式化验证码Ticket的Redis键名
// 格式: <服务名>:captcha:ticket:<key>
// 示例: edu_mall:captcha:ticket:abc123
func fmtVerifyCaptchaTicket(key string) string {
	return fmt.Sprintf("%s:captcha:ticket:%s", config.ServerFullName, key)
}

// SetCaptchaKey 存储验证码Key到Redis
// 参数:
//   - ctx: 上下文
//   - key: 验证码标识
//   - value: 验证码答案(JSON格式)
//   - expire: 过期时间
// 返回: 错误信息
// 用途: 存储滑块验证码的正确答案
func (v *Verify) SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaKey(key)
	return v.redis.Set(redisKey, value, expire).Err()
}

// GetCaptchaKey 获取验证码Key并删除
// 参数:
//   - ctx: 上下文
//   - key: 验证码标识
// 返回: 验证码答案(JSON格式)和错误信息
// 特性: 获取后立即删除,防止重复使用
// 调用链: service.CheckCaptcha -> GetCaptchaKey
func (v *Verify) GetCaptchaKey(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	// 获取后删除,确保验证码只能使用一次
	v.redis.Del(redisKey)
	return get, nil
}

// SetCaptchaTicket 存储验证码Ticket到Redis
// 参数:
//   - ctx: 上下文
//   - key: Ticket标识
//   - value: Ticket内容
//   - expire: 过期时间
// 返回: 错误信息
// 用途: 验证通过后生成临时凭证,用于后续登录
func (v *Verify) SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaTicket(key)
	return v.redis.Set(redisKey, value, expire).Err()
}

// GetCaptchaTicket 获取验证码Ticket并删除
// 参数:
//   - ctx: 上下文
//   - key: Ticket标识
// 返回: Ticket内容和错误信息
// 特性: 获取后立即删除,防止重复使用
// 调用链: service.Login -> GetCaptchaTicket
func (v *Verify) GetCaptchaTicket(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaTicket(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	// 获取后删除,确保Ticket只能使用一次
	v.redis.Del(redisKey)
	return get, nil
}
