// Package config 配置管理模块
// 职责: 统一管理应用配置,支持本地YAML文件和etcd远程配置
// 特性: 支持etcd配置热更新
package config

import (
	"flag"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"time"
)

const (
	ServerName     = "mall"         // 服务名称
	ServerFullName = "edu.mall"     // 服务完整名称,用于etcd路径
)

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", ServerFullName) // etcd配置键
	etcdAddr        string                                               // etcd地址,通过命令行参数-r或环境变量ETCD_ADDR指定
	localConfigPath string                                               // 本地配置文件路径,默认mall_local.yml
	GlobalConfig    Config                                               // 全局配置对象,用于热更新
)

// Config 应用配置结构体
type Config struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
}

// Server HTTP服务器配置
type Server struct {
	HttpPort    int    `yaml:"http_port"`    // HTTP服务端口
	Env         string `yaml:"env"`          // 环境标识: dev/test/prod
	EnablePprof bool   `yaml:"enable_pprof"` // 是否启用pprof性能分析
	LogLevel    string `yaml:"log_level"`    // 日志级别: debug/info/warn/error
}

// Mysql 数据库配置
type Mysql struct {
	Dialect  string `yaml:"dialect"`  // 数据库类型,默认mysql
	User     string `yaml:"user"`     // 用户名
	Password string `yaml:"password"` // 密码
	Host     string `yaml:"host"`     // 主机地址
	Port     int    `yaml:"port"`     // 端口
	Database string `yaml:"database"` // 数据库名
	Charset  string `yaml:"charset"`  // 字符集
	ShowSql  bool   `yaml:"show_sql"` // 是否打印SQL
	MaxOpen  int    `yaml:"max_open"` // 最大打开连接数
	MaxIdle  int    `yaml:"max_idle"` // 最大空闲连接数
}

// GetDsn 生成MySQL DSN连接字符串
// 格式: user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=true&loc=Local
func (m *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

// Redis 缓存配置
type Redis struct {
	Addr    string `yaml:"addr"`     // Redis地址,格式: host:port
	PWD     string `yaml:"password"` // 密码
	DBIndex int    `yaml:"db_index"` // 数据库索引
	MaxIdle int    `yaml:"max_idle"` // 最大空闲连接数
	MaxOpen int    `yaml:"max_open"` // 最大活跃连接数
}

// init 初始化命令行参数
// -c: 指定本地配置文件路径,默认mall_local.yml
// -r: 指定etcd地址,默认从环境变量ETCD_ADDR获取
func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default consul address")
}

// InitConfig 初始化配置
// 优先级: etcd远程配置 > 本地YAML文件
// 返回: 配置对象指针
// 调用: main.main -> InitConfig
func InitConfig() *Config {
	var (
		err      error
		tempConf = &Config{}
		vipConf  = viper.New()
	)

	flag.Parse()

	// 如果指定了etcd地址,优先使用etcd配置
	if etcdAddr != "" {
		tempConf, err = getFromRemoteAndWatchUpdate(vipConf)
		if err != nil {
			panic(err)
		}
		return tempConf
	}

	// 否则从本地文件加载
	tempConf, err = getFromLocal()
	if err != nil {
		panic(err)
	}
	return tempConf
}

// getFromRemoteAndWatchUpdate 从etcd获取配置并监听更新
// 参数: v viper实例
// 返回: 配置对象和错误
// 特性: 启动协程每分钟检查一次配置更新,自动热更新GlobalConfig
func getFromRemoteAndWatchUpdate(v *viper.Viper) (*Config, error) {
	tempConf := Config{}
	if err := v.AddRemoteProvider("etcd3", etcdAddr, etcdKey); err != nil {
		return nil, err
	}
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}

	// 反序列化配置到结构体
	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}

	// 启动协程监听配置变更,实现热更新
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			if err := v.WatchRemoteConfig(); err == nil {
				_ = v.Unmarshal(&GlobalConfig)
				fmt.Println(">>> etcd config hot-reloaded: ", gconv.String(GlobalConfig))
			}
		}
	}()
	return &tempConf, nil
}

// getFromLocal 从本地YAML文件加载配置
// 返回: 配置对象和错误
// 文件路径由命令行参数-c指定,默认mall_local.yml
func getFromLocal() (*Config, error) {
	tempConf := Config{}
	if _, err := os.Stat(localConfigPath); err == nil {
		content, err := os.ReadFile(localConfigPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(content, &tempConf)
		if err != nil {
			return nil, err
		}
		return &tempConf, nil
	}
	return nil, fmt.Errorf("local config file not found ,file_name: %s", localConfigPath)
}
