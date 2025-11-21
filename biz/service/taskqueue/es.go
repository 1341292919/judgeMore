package taskqueue

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"judgeMore/biz/dal/es"
	"judgeMore/biz/dal/mysql"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/taskqueue"
)

// 异步同步es数据库 主要用于学校认定的奖项
func AddUpdateElasticTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateElastic(ctx)
	}})
}
func updateElastic(ctx context.Context) error {
	reList, _, err := mysql.QueryRecognizedEvent(ctx)
	if err != nil {
		logger.Infof("taskqueue es failed : checked mysql %v", err.Error())
		return err
	}
	for _, r := range reList {
		if !r.IsActive {
			err = es.RemoveItem(ctx, constants.IndexName, r.RecognizedEventId)
			if err != nil {
				logger.Infof("taskqueue es failed :remove es %v", err.Error())
				return err
			}
			continue
		}
		err = es.AddItem(ctx, constants.IndexName, r)
		if err != nil {
			logger.Infof("taskqueue es failed : insert es %v", err.Error())
			return err
		}
	}
	return nil
}
