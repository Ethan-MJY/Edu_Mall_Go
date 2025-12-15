// Package admin 管理员数据访问层
// 职责: 封装admin_user表的CRUD操作
// 调用链: service -> repo -> GORM
package admin

import (
	"context"
	"mall/adaptor"
	"mall/adaptor/repo/model"
	"mall/adaptor/repo/query"
	"mall/consts"
	"mall/service/do"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// IAdminUser 管理员用户数据访问接口
type IAdminUser interface {
	CreateUser(ctx context.Context, req *do.CreateUser) (int64, error)        // 创建管理员
	UpdateUser(ctx context.Context, req *do.UpdateUser) error                 // 更新管理员信息
	UpdateUserStatus(ctx context.Context, req *do.UpdateUserStatus) error     // 更新管理员状态(启用/禁用)
	UpdateUserPassword(ctx context.Context, req *do.UpdateUserPassword) error // 更新管理员密码
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)  // 获取管理员详细信息
}

// AdminUser 管理员用户数据访问实现
type AdminUser struct {
	db    *gorm.DB      // 数据库连接
	redis *redis.Client // Redis客户端(预留用于缓存)
}

// NewAdminUser 创建管理员用户数据访问实例
// 参数: adaptor 适配器,提供数据库和Redis连接
// 返回: AdminUser实例
// 调用链: service.NewService -> NewAdminUser
func NewAdminUser(adaptor adaptor.IAdaptor) *AdminUser {
	return &AdminUser{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
}

// CreateUser 创建管理员用户
// 参数:
//   - ctx: 上下文
//   - req: 创建用户请求DO对象
//
// 返回: 用户ID和错误信息
// 业务逻辑:
//  1. 记录创建时间和更新时间
//  2. 默认状态为启用
//  3. 记录创建人和更新人
//
// 调用链: service.CreateUser -> repo.CreateUser -> GORM.Create
func (a *AdminUser) CreateUser(ctx context.Context, req *do.CreateUser) (int64, error) {
	timeNow := time.Now()
	qs := query.Use(a.db).AdminUser
	addObj := &model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		UpdateBy: req.AdminUserID, // 记录创建人
		Status:   consts.IsEnable, // 默认启用
		CreateBy: req.AdminUserID,
	}
	err := qs.WithContext(ctx).Create(addObj)
	if err != nil {
		return 0, err
	}
	return addObj.ID, nil
}

// UpdateUser 更新管理员用户基本信息
// 参数:
//   - ctx: 上下文
//   - req: 更新用户请求DO对象
//
// 返回: 错误信息
// 可更新字段: 姓名、昵称、性别
// 自动更新: 更新时间、更新人
// 调用链: service.UpdateUser -> repo.UpdateUser -> GORM.Updates
func (a *AdminUser) UpdateUser(ctx context.Context, req *do.UpdateUser) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Sex:      req.Sex,
		UpdateAt: time.Now(),
		UpdateBy: req.AdminUserID, // 记录更新人
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserStatus 更新管理员用户状态
// 参数:
//   - ctx: 上下文
//   - req: 更新状态请求DO对象
//
// 返回: 错误信息
// 状态值: consts.IsEnable(1)启用, consts.IsDisable(-1)禁用
// 用途: 管理员账号的启用/停用管理
// 调用链: service.UpdateUserStatus -> repo.UpdateUserStatus -> GORM.Updates
func (a *AdminUser) UpdateUserStatus(ctx context.Context, req *do.UpdateUserStatus) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(model.AdminUser{
		Status:   req.Status,
		UpdateAt: time.Now(),
		UpdateBy: req.AdminUserID, // 记录更新人
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserPassword 更新管理员用户密码
// 参数:
//   - ctx: 上下文
//   - req: 更新密码请求DO对象
//
// 返回: 错误信息
// 注意: 传入的password应该已经是SHA256哈希后的值
// 调用链: service.ResetPassword -> repo.UpdateUserPassword -> GORM.Updates
func (a *AdminUser) UpdateUserPassword(ctx context.Context, req *do.UpdateUserPassword) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(model.AdminUser{
		Password: req.Password, // 哈希后的密码
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserInfo 获取管理员用户详细信息
// 参数:
//   - ctx: 上下文
//   - userId: 用户ID
//
// 返回: 用户对象和错误信息
// 用途: 获取管理员个人资料、权限查询等
// 调用链: service.GetUserInfo -> repo.GetUserInfo -> GORM.First
func (a *AdminUser) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}
