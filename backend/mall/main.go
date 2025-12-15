// Package main 应用程序入口
// 职责: 初始化配置、数据库连接、Redis连接,启动HTTP服务器
package main

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mall/adaptor"
	"mall/config"
	"mall/router"
	"mall/utils/logger"
)

// main 应用程序主入口
// 执行流程:
// 1. 初始化配置(支持本地文件和etcd)
// 2. 设置日志级别
// 3. 初始化MySQL连接
// 4. 初始化Redis连接
// 5. 启动HTTP服务器
func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.Server.LogLevel)

	dbClient, err := initMysql(&conf.Mysql)
	handleErr(err)
	logger.Debug("mysql connect success")

	rdsClient, err := initRedis(&conf.Redis)
	handleErr(err)
	logger.Debug("client connect success")

	startServer(conf, dbClient, rdsClient).Run()
}

// startServer 启动HTTP服务器
// 参数:
//   - conf: 配置对象
//   - db: GORM数据库连接
//   - redis: Redis客户端
//
// 返回: router.App HTTP服务器实例
// 调用链: main -> router.NewApp -> router.NewRouter -> adaptor.NewAdaptor
func startServer(conf *config.Config, db *gorm.DB, redis *redis.Client) *router.App {
	return router.NewApp(conf.Server.HttpPort,
		router.NewRouter(
			conf,
			adaptor.NewAdaptor(conf, db, redis),
			// 健康检查函数: 用于/ping接口检测MySQL和Redis连通性
			func() error {
				err := func() error {
					pingDb, err := db.DB()
					handleErr(err)
					return pingDb.Ping()
				}()
				if err != nil {
					return errors.New("mysql connect failed")
				}
				return redis.Ping().Err()
			},
		),
	)
}

// initRedis 初始化Redis连接
// 参数: conf Redis配置
// 返回: Redis客户端实例和错误
// 连接失败时返回error
func initRedis(conf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.PWD,
		DB:           conf.DBIndex,
		MinIdleConns: conf.MaxIdle,
		PoolSize:     conf.MaxOpen,
	})
	// 通过PING命令验证连接
	if r, _ := client.Ping().Result(); r != "PONG" {
		return nil, errors.New("redis connect failed")
	}
	return client, nil
}

// initMysql 初始化MySQL连接
// 参数: conf MySQL配置
// 返回: GORM数据库实例和错误
// 连接池配置:
//   - MaxIdle: 最小值5,配置值+1
//   - MaxOpen: 最小值10,配置值+1
func initMysql(conf *config.Mysql) (*gorm.DB, error) {
	// 确保连接池配置合理的最小值
	conf.MaxIdle = lo.Max([]int{conf.MaxIdle + 1, 5})
	conf.MaxOpen = lo.Max([]int{conf.MaxOpen + 1, 10})
	dsn := conf.GetDsn()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 通过Ping验证连接
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetMaxOpenConns(conf.MaxOpen)
	return db, nil
}

// handleErr 错误处理函数
// 如果错误不为nil,则panic中止程序
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
