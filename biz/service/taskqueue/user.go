package taskqueue

import (
	"context"
	"fmt"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/taskqueue"
)

func AddStuToAdminRelation(ctx context.Context, key string, u *model.User) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateStuToAdminRelation(ctx, u)
	}})
}

func updateStuToAdminRelation(ctx context.Context, u *model.User) error {
	relationList, _, err := mysql.QueryAllRelation(ctx)
	if err != nil {
		return err
	}
	// 专业与年级
	var c []string
	for _, r := range relationList {
		if r.CollegeName == u.College {
			c = append(c, r.UserId)
			continue
		}
		if r.MajorName == u.Major && r.Grade == u.Grade {
			c = append(c, r.UserId)
		}
	}
	if len(c) == 0 {
		return nil
	}
	err = mysql.InsertStuIdToAdminRelation(ctx, u.Uid, c)
	if err != nil {
		fmt.Printf("updateStuToAdminRelation failed : %v", err)
		return err
	}
	return nil
}
