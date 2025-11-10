package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
)

type ScoreService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewScoreService(ctx context.Context, c *app.RequestContext) *ScoreService {
	return &ScoreService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *ScoreService) QueryScoreRecordByScoreId(score_id string) (*model.ScoreRecord, error) {
	exist, err := mysql.IsScoreRecordExist(svc.ctx, score_id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceRecordNotExistCode, "Socre Result not exist")
	}
	recordInfo, err := mysql.QueryScoreRecordByScoreId(svc.ctx, score_id)
	if err != nil {
		return nil, err
	}
	return recordInfo, nil
}

func (svc *ScoreService) QueryScoreRecordByEventId(event_id string) (*model.ScoreRecord, error) {
	exist, err := mysql.IsScoreRecordExist_Event(svc.ctx, event_id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceRecordNotExistCode, "Socre Result not exist")
	}
	recordInfo, err := mysql.QueryScoreRecordByEventId(svc.ctx, event_id)
	if err != nil {
		return nil, err
	}
	return recordInfo, nil
}

func (svc *ScoreService) QueryScoreRecordByStuId(stu_id string) ([]*model.ScoreRecord, int64, error) {
	exist, err := mysql.IsUserExist(svc.ctx, &model.User{Uid: stu_id})
	if err != nil {
		return nil, -1, err
	}
	if !exist {
		return nil, -1, errno.NewErrNo(errno.ServiceUserExistCode, "user not exist")
	}
	recordInfoList, count, err := mysql.QueryScoreRecordByStuId(svc.ctx, stu_id)
	if err != nil {
		return nil, -1, err
	}
	return recordInfoList, count, nil
}

// 用于直接修改积分。
func (svc *ScoreService) ReviseScore(result_id string, score float64) error {
	exist, err := mysql.IsScoreRecordExist(svc.ctx, result_id)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceRecordNotExistCode, "Socre Result not exist")
	}
	info, err := mysql.QueryScoreRecordByScoreId(svc.ctx, result_id)
	if err != nil {
		return err
	}
	// 验证用户权限
	user_id := GetUserIDFromContext(svc.c)
	exist, err = mysql.IsAdminRelationExist(svc.ctx, user_id, info.UserId)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceNoAuthToDo, "No permission to update the stu's score")
	}
	err = mysql.UpdatesScore(svc.ctx, result_id, score)
	if err != nil {
		return err
	}
	return nil
}
