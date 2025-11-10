package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
	"time"
)

// 学校结构以及权责结构也是高频访问 cache存储

func IsMajorExist(ctx context.Context) (bool, error) {
	keys, err := structureCa.Keys(ctx, "major_*").Result()
	if err != nil {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "get rule keys from redis error:"+err.Error())
	}
	if len(keys) == 0 {
		return false, nil
	}
	return true, nil
}

func MajorToCache(ctx context.Context, majorList []*model.Major) error {
	for _, r := range majorList {
		key := fmt.Sprintf("major_%v", r.MajorId)
		// 使用 JSON 序列化
		info, err := json.Marshal(r)
		if err != nil {
			return errno.NewErrNo(errno.InternalServiceErrorCode, "marshal major to json error:"+err.Error())
		}
		expiration := 72 * time.Hour
		err = structureCa.Set(ctx, key, info, expiration).Err()
		if err != nil {
			return errno.NewErrNo(errno.InternalRedisErrorCode, "write major to cache error:"+err.Error())
		}
	}
	return nil
}

func QueryAllMajor(ctx context.Context) ([]*model.Major, error) {
	keys, err := structureCa.Keys(ctx, "major_*").Result()
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:get major key error:"+err.Error())
	}
	pipe := structureCa.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query major fail:"+err.Error())
	}
	majors := make([]*model.Major, 0)
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		data, err := getCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			continue
		}
		var major model.Major
		err = json.Unmarshal([]byte(data), &major)
		if err != nil {
			continue
		}
		majors = append(majors, &major)
	}
	return majors, nil
}

// 学院

func IsCollegeExist(ctx context.Context) (bool, error) {
	keys, err := structureCa.Keys(ctx, "college_*").Result()
	if err != nil {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "get rule keys from redis error:"+err.Error())
	}
	if len(keys) == 0 {
		return false, nil
	}
	return true, nil
}

func CollegeToCache(ctx context.Context, collegeList []*model.College) error {
	for _, r := range collegeList {
		key := fmt.Sprintf("college_%v", r.CollegeId)
		// 使用 JSON 序列化
		info, err := json.Marshal(r)
		if err != nil {
			return errno.NewErrNo(errno.InternalServiceErrorCode, "marshal college to json error:"+err.Error())
		}
		expiration := 72 * time.Hour
		err = structureCa.Set(ctx, key, info, expiration).Err()
		if err != nil {
			return errno.NewErrNo(errno.InternalRedisErrorCode, "write college to cache error:"+err.Error())
		}
	}
	return nil
}

func QueryAllCollege(ctx context.Context) ([]*model.College, error) {
	keys, err := structureCa.Keys(ctx, "college_*").Result()
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:get college key error:"+err.Error())
	}
	pipe := structureCa.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query college fail:"+err.Error())
	}
	colleges := make([]*model.College, 0)
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		data, err := getCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			continue
		}
		var college model.College
		err = json.Unmarshal([]byte(data), &college)
		if err != nil {
			continue
		}
		colleges = append(colleges, &college)
	}
	return colleges, nil
}

// relation
func IsRelationExist(ctx context.Context) (bool, error) {
	keys, err := structureCa.Keys(ctx, "relation_*").Result()
	if err != nil {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "get rule keys from redis error:"+err.Error())
	}
	if len(keys) == 0 {
		return false, nil
	}
	return true, nil
}

func RelationToCache(ctx context.Context, relationList []*model.Relation) error {
	for _, r := range relationList {
		key := fmt.Sprintf("relation_%v_%v", r.UserId, r.RelationId)
		// 使用 JSON 序列化
		info, err := json.Marshal(r)
		if err != nil {
			return errno.NewErrNo(errno.InternalServiceErrorCode, "marshal relation to json error:"+err.Error())
		}
		expiration := 72 * time.Hour
		err = structureCa.Set(ctx, key, info, expiration).Err()
		if err != nil {
			return errno.NewErrNo(errno.InternalRedisErrorCode, "write relation to cache error:"+err.Error())
		}
	}
	return nil
}

func QueryAllRelation(ctx context.Context) ([]*model.Relation, error) {
	keys, err := structureCa.Keys(ctx, "relation_*").Result()
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:get relation key error:"+err.Error())
	}
	pipe := structureCa.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query relation fail:"+err.Error())
	}
	relations := make([]*model.Relation, 0)
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		data, err := getCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			continue
		}
		var relation model.Relation
		err = json.Unmarshal([]byte(data), &relation)
		if err != nil {
			continue
		}
		relations = append(relations, &relation)
	}
	return relations, nil
}

func QueryRelationById(ctx context.Context, user_id string) ([]*model.Relation, error) {
	pattern := fmt.Sprintf("relation_%v_*", user_id)
	keys, err := structureCa.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:get relation key error:"+err.Error())
	}
	pipe := structureCa.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query relation fail:"+err.Error())
	}
	relations := make([]*model.Relation, 0)
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		data, err := getCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			continue
		}
		var relation model.Relation
		err = json.Unmarshal([]byte(data), &relation)
		if err != nil {
			continue
		}
		relations = append(relations, &relation)
	}
	return relations, nil
}
