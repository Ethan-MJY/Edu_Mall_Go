// Package tools 通用工具函数模块
// 职责: 提供UUID生成等通用工具函数
package tools

import (
	"github.com/google/uuid"
	"strings"
)

// UUIDHex 生成32位十六进制UUID
// 返回: 去除短横线的UUID字符串,32位十六进制
// 用途: 文件key、订单号等唯一标识生成
func UUIDHex() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
