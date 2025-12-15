// Package admin 管理员业务逻辑层-用户管理
// 职责: 管理员用户的CRUD业务逻辑
package admin

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/common"
	"mall/service/do"
	"mall/service/dto"
	"mall/utils/logger"
)

// CreateUser 创建管理员用户
// 参数:
//   - ctx: 上下文
//   - adminUser: 当前操作的管理员(从Token解析)
//   - req: 创建用户请求DTO
// 返回: 用户ID和错误码
// 业务流程:
//   1. 转换DTO为DO对象
//   2. 记录操作人ID
//   3. 调用数据访问层创建用户
// 调用链: api.CreateUser -> service.CreateUser -> repo.CreateUser
func (s *Service) CreateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.CreateUserReq) (int64, common.Errno) {
	userID, err := s.adminUser.CreateUser(ctx, &do.CreateUser{
		AdminUserID: adminUser.UserID, // 记录创建人ID
		Name:        req.Name,
		NickName:    req.NickName,
		Mobile:      req.Mobile,
		Sex:         req.Sex,
	})
	if err != nil {
		logger.Error("CreateUser error", zap.Error(err), zap.Any("req", req))
		return 0, common.DatabaseErr.WithErr(err)
	}
	return userID, common.OK
}

// UpdateUser 更新管理员用户基本信息
// 参数:
//   - ctx: 上下文
//   - adminUser: 当前操作的管理员
//   - req: 更新用户请求DTO
// 返回: 错误码
// 可更新字段: 姓名、昵称、性别
// 调用链: api.UpdateUser -> service.UpdateUser -> repo.UpdateUser
func (s *Service) UpdateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserReq) common.Errno {
	err := s.adminUser.UpdateUser(ctx, &do.UpdateUser{
		ID:          req.ID,
		Name:        req.Name,
		NickName:    req.NickName,
		Sex:         req.Sex,
		AdminUserID: adminUser.UserID, // 记录更新人ID
	})
	if err != nil {
		logger.Error("UpdateUser error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

// UpdateUserStatus 更新管理员用户状态
// 参数:
//   - ctx: 上下文
//   - adminUser: 当前操作的管理员
//   - req: 更新状态请求DTO
// 返回: 错误码
// 状态值: 1(启用) / -1(禁用)
// 调用链: api.UpdateUserStatus -> service.UpdateUserStatus -> repo.UpdateUserStatus
func (s *Service) UpdateUserStatus(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserStatusReq) common.Errno {
	err := s.adminUser.UpdateUserStatus(ctx, &do.UpdateUserStatus{
		ID:          req.ID,
		Status:      req.Status,
		AdminUserID: adminUser.UserID, // 记录更新人ID
	})
	if err != nil {
		logger.Error("UpdateUserStatus error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

// GetUserInfo 获取管理员用户详细信息
// 参数:
//   - ctx: 上下文
//   - adminUser: 当前登录的管理员
// 返回: 用户信息DTO和错误码
// TODO: 当前写死查询ID=1,应改为查询当前用户
// 调用链: api.GetUserInfo -> service.GetUserInfo -> repo.GetUserInfo
func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser) (*dto.UserInfoResp, common.Errno) {
	// TODO: 应该查询adminUser.UserID,而不是写死1
	user, err := s.adminUser.GetUserInfo(ctx, 1)
	if err != nil {
		// 用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo GetUserInfo error", zap.Error(err), zap.Any("user_id", adminUser))
		return nil, common.DatabaseErr.WithErr(err)
	}
	return &dto.UserInfoResp{
		Name:   user.Name,
		UserID: user.ID,
	}, common.OK
}
