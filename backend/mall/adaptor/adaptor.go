// Package adaptor 适配器模块
// 职责: 统一管理外部依赖(数据库、Redis、配置)的访问接口
// 设计模式: 适配器模式,隔离业务层与基础设施层,实现依赖注入
package adaptor

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"mall/config"
)

// IAdaptor 适配器接口
// 提供统一的访问入口,供上层(service/api)获取基础设施依赖
type IAdaptor interface {
	GetConfig() *config.Config // 获取配置对象
	GetDB() *gorm.DB           // 获取数据库连接
	GetRedis() *redis.Client   // 获取Redis客户端
}

// Adaptor 适配器实现
// 持有配置、数据库、Redis三大基础设施对象
type Adaptor struct {
	conf  *config.Config   // 配置对象
	db    *gorm.DB         // 数据库连接(GORM)
	redis *redis.Client    // Redis客户端
}

// NewAdaptor 创建适配器实例
// 参数:
//   - conf: 配置对象
//   - db: GORM数据库连接
//   - redis: Redis客户端
// 返回: Adaptor实例
// 调用链: main.main -> NewAdaptor
func NewAdaptor(conf *config.Config, db *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		conf:  conf,
		db:    db,
		redis: redis,
	}
}

// GetConfig 获取配置对象
// 返回: 配置对象指针
func (a *Adaptor) GetConfig() *config.Config {
	return a.conf
}

// GetDB 获取数据库连接
// 返回: GORM数据库连接对象
func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}

// GetRedis 获取Redis客户端
// 返回: Redis客户端对象
func (a *Adaptor) GetRedis() *redis.Client {
	return a.redis
}
