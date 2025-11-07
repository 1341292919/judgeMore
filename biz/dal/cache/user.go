package cache

import (
	"context"
	"fmt"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/utils"
	"strings"
	"time"
)

func IsKeyExist(ctx context.Context, key string) bool {
	return userCa.Exists(ctx, key).Val() == 1
}

func GetCodeCache(ctx context.Context, key string) (code string, err error) {
	value, err := userCa.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "write code to cache error:"+err.Error())
	}
	var storedCode string
	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid code format, expected 2 parts, got %d", len(parts))
	}
	storedCode = parts[0]
	return storedCode, nil
}
func PutCodeToCache(ctx context.Context, key string) (code string, err error) {
	code = utils.GenerateRandomCode(6)
	timeNow := time.Now().Unix()
	value := fmt.Sprintf("%s_%d", code, timeNow)
	expiration := 2 * time.Minute
	err = userCa.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "write code to cache error:"+err.Error())
	}
	return code, nil
}

func DeleteCodeCache(ctx context.Context, key string) error {
	err := userCa.Del(ctx, key).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "delete code from cache error:"+err.Error())
	}
	return nil
}

func PutTokenIdToCache(ctx context.Context, key string) error {
	value := fmt.Sprintf("%s", time.Now().Unix())
	expiration := 72 * time.Hour // 与refresh-token过期时间一致
	err := userCa.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, "write token id to cache error:"+err.Error())
	}
	return nil
}
