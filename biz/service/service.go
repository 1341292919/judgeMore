package service

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/pkg/constants"
)

func GetUserIDFromContext(c *app.RequestContext) string {
	if c == nil || c.Keys == nil {
		return ""
	}

	data, exists := c.Keys[constants.ContextUserId]
	if !exists {
		return ""
	}

	// 类型断言确保返回的是 string
	if userID, ok := data.(string); ok {
		return userID
	}

	return ""
}
func GetTokenIdFromContext(c *app.RequestContext) string {
	if c == nil || c.Keys == nil {
		return ""
	}
	data, exists := c.Keys[constants.ContextTokenId]

	if !exists {
		return ""
	}
	if tokenID, ok := data.(string); ok {
		return tokenID
	}

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
