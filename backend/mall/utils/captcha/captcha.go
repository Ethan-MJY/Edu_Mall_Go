// Package captcha 验证码工具模块
// 职责: 初始化和配置滑块验证码
package captcha

import (
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

// NewSlideCaptcha 创建滑块验证码实例
// 返回: 配置好的滑块验证码对象
// 配置:
//   - 单图模式(GenGraphNumber=1)
//   - 使用内置背景图片和滑块图形
//
// 调用: service/admin/service.NewService -> NewSlideCaptcha
func NewSlideCaptcha() slide.Captcha {
	builder := slide.NewBuilder(
		slide.WithGenGraphNumber(1), // 生成1个滑块图形
	)

	// 加载背景图片资源
	imgs, err := imagesv2.GetImages()
	if err != nil {
		panic(err)
	}

	// 加载滑块图形资源
	graphs, err := tiles.GetTiles()
	if err != nil {
		panic(err)
	}

	// 转换图形格式
	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for _, g := range graphs {
		newGraphs = append(newGraphs, &slide.GraphImage{
			MaskImage:    g.MaskImage,
			OverlayImage: g.OverlayImage,
			ShadowImage:  g.ShadowImage,
		})
	}

	// 设置资源并构建
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)
	return builder.Make()
}
