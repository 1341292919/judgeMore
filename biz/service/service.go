package service

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/pkg/constants"
)

// 这里其实不允许返回为空值，返回空值必然导致后续业务错误
func GetUserIDFromContext(c *app.RequestContext) string {
	if c == nil || c.Keys == nil {
		panic(fmt.Errorf("stream c or c.key is nil"))
		return ""
	}

	data, exists := c.Keys[constants.ContextUserId]
	if !exists {
		panic(fmt.Errorf("userId is nil"))
		return ""
	}

	// 类型断言确保返回的是 string
	if userID, ok := data.(string); ok {
		return userID
	}

	panic(fmt.Errorf("userId is not string"))
	return ""
}
func GetTokenIdFromContext(c *app.RequestContext) string {
	if c == nil || c.Keys == nil {
		panic(fmt.Errorf("stream c or c.key is nil"))
		return ""
	}
	data, exists := c.Keys[constants.ContextTokenId]

	if !exists {
		panic(fmt.Errorf("tokenId is nil"))
		return ""
	}
	if tokenID, ok := data.(string); ok {
		return tokenID
	}
	panic(fmt.Errorf("tokenId is not string"))
	return ""
}
func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("无法转换为int64，类型为 %T", value)
	}
}
