package utils

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"judgeMore/config"
	"math/rand"
	"net"
	"strings"
	"time"
)

// GetMysqlDSN 会拼接 mysql 的 DSN
func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	return dsn, nil
}

// AddrCheck 会检查当前的监听地址是否已被占用
func AddrCheck(addr string) bool {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.Errorf("utils.AddrCheck: failed to close listener: %v", err.Error())
		}
	}()
	return true
}

// GetAvailablePort 会尝试获取可用的监听地址
func GetAvailablePort() (string, error) {
	if config.Service.AddrList == nil {
		return "", errors.New("utils.GetAvailablePort: config.Service.AddrList is nil")
	}
	for _, addr := range config.Service.AddrList {
		if ok := AddrCheck(addr); ok {
			return addr, nil
		}
	}
	return "", errors.New("utils.GetAvailablePort: not available port from config")
}

var m = map[string]int{
	"gateway":   0,
	"video":     1,
	"user":      2,
	"interact":  3,
	"websocket": 4,
}

// 生成指定位数的随机验证码（字母+数字）
func GenerateRandomCode(length int) string {
	// 字符集：26个小写字母 + 26个大写字母 + 10个数字
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charSet[r.Intn(len(charSet))]
	}

	return string(code)
}

func VerifyPageParam(page_num, page_size int64) bool {
	if page_num < 0 || page_size < 0 {
		return false
	}
	return true
}
