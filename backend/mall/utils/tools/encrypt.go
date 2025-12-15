// Package tools 加密工具模块
// 职责: 提供各类加密算法实现
package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hash SHA256哈希计算
// 参数: text 待哈希的明文字符串
// 返回: 64位十六进制哈希字符串
// 用途: 手机号加密存储、密码哈希等
func Sha256Hash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
