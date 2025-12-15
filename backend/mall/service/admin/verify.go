// Package admin 管理员业务逻辑层-验证码
// 职责: 滑块验证码的生成和校验业务逻辑
package admin

import (
	"context"
	"encoding/json"
	"github.com/wenlng/go-captcha/v2/slide"
	"go.uber.org/zap"
	"mall/common"
	"mall/service/dto"
	"mall/utils/logger"
	"mall/utils/tools"
	"time"
)

// GetSlideCaptcha 获取滑块验证码
// 参数: ctx 上下文
// 返回: 验证码响应DTO和错误码
// 业务流程:
//   1. 生成滑块验证码(背景图+滑块图)
//   2. 获取滑块正确位置坐标
//   3. 将坐标JSON序列化后存入Redis(key为UUID,有效期2分钟)
//   4. 返回验证码图片Base64和滑块尺寸信息
// 调用链: api.GetSmsCodeCaptcha -> service.GetSlideCaptcha
func (s *Service) GetSlideCaptcha(ctx context.Context) (*dto.GetVerifyCaptchaResp, common.Errno) {
	// 1. 生成验证码
	captData, err := s.captcha.Generate()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}

	// 2. 获取滑块正确位置数据
	dotData := captData.GetData()
	if dotData == nil {
		logger.Error("GetSlideCaptcha GetData error")
		return nil, common.ServerErr.WithMsg("GetData is nil")
	}

	// 3. 将坐标数据序列化为JSON
	dots, err := json.Marshal(dotData)
	if err != nil {
		logger.Error("GetSlideCaptcha json.Marshal error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}

	// 4. 获取背景图和滑块图的Base64编码
	var mBs64Data, tBs64Data string
	mBs64Data, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetMasterImage error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}
	tBs64Data, err = captData.GetTileImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetTileImage error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}

	// 5. 生成唯一Key并存入Redis
	key := tools.UUIDHex()
	err = s.verify.SetCaptchaKey(ctx, key, string(dots), time.Minute*2) // 有效期2分钟
	if err != nil {
		logger.Error("GetSlideCaptcha SetCaptchaKey error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}

	// 6. 返回验证码数据
	return &dto.GetVerifyCaptchaResp{
		Key:            key,          // 验证码唯一标识
		ImageBs64:      mBs64Data,    // 背景图Base64
		TitleImageBs64: tBs64Data,    // 滑块图Base64
		TitleHeight:    dotData.Height, // 滑块高度
		TitleWidth:     dotData.Width,  // 滑块宽度
		TitleX:         dotData.TileX,  // 滑块初始X坐标
		TitleY:         dotData.TileY,  // 滑块初始Y坐标
		Expire:         110,            // 前端显示的剩余秒数
	}, common.OK
}

// CheckSlideCaptcha 校验滑块验证码
// 参数:
//   - ctx: 上下文
//   - req: 校验请求DTO(包含key和用户滑动的坐标)
// 返回: 校验响应DTO和错误码
// 业务流程:
//   1. 从Redis获取正确坐标(获取后自动删除)
//   2. 反序列化坐标数据
//   3. 校验用户滑动坐标与正确坐标的误差(允许5像素误差)
//   4. 校验成功生成Ticket存入Redis(有效期5分钟)
//   5. 返回Ticket用于后续登录
// 调用链: api.CheckSmsCodeCaptcha -> service.CheckSlideCaptcha
func (s *Service) CheckSlideCaptcha(ctx context.Context, req *dto.CheckCaptchaReq) (*dto.CheckCaptchaDtoResp, common.Errno) {
	// 1. 从Redis获取验证码正确坐标(获取后自动删除)
	captData, err := s.verify.GetCaptchaKey(ctx, req.Key)
	if err != nil {
		logger.Error("CheckSlideCaptcha GetCaptchaKey error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}
	if captData == "" {
		return nil, common.ParamErr.WithMsg("滑块已过期，请刷新重试")
	}

	// 2. 反序列化坐标数据
	dot := slide.Block{}
	err = json.Unmarshal([]byte(captData), &dot)
	if err != nil {
		logger.Error("CheckSlideCaptcha json.Unmarshal error", zap.Error(err))
		return nil, common.InvalidCaptchaErr
	}

	// 3. 校验坐标(允许5像素误差)
	ok := slide.CheckPoint(int64(req.SlideX), int64(req.SlideY), int64(dot.X), int64(dot.Y), 5)
	if !ok {
		return nil, common.InvalidCaptchaErr
	}

	// 4. 生成Ticket并存入Redis(有效期5分钟)
	ticket := tools.UUIDHex()
	err = s.verify.SetCaptchaTicket(ctx, ticket, req.Key, time.Minute*5)
	if err != nil {
		logger.Error("CheckSlideCaptcha SetCaptchaTicket error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}

	// 5. 返回Ticket
	return &dto.CheckCaptchaDtoResp{
		Ticket: ticket, // 验证通过凭证,用于登录
		Expire: 280,    // 前端显示的剩余秒数
	}, common.OK
}
