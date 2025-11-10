package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/biz/service/taskqueue"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

type AppealService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewAppealService(ctx context.Context, c *app.RequestContext) *AppealService {
	return &AppealService{
		ctx: ctx,
		c:   c,
	}
}
func (svc *AppealService) NewAppeal(a *model.Appeal) (string, error) {
	// 首先检查该记录是否已经进行申诉
	exist, err := mysql.IsAppealExist(svc.ctx, a.ResultId)
	if err != nil {
		return "", err
	}
	if exist {
		return "", errno.NewErrNo(errno.ServiceAppealExistCode, "result already appeal")
	}
	// 这里应该完成一次验证 验证result属于该user
	stu_id := GetUserIDFromContext(svc.c)
	a.UserId = stu_id
	resultInfo, err := mysql.QueryScoreRecordByScoreId(svc.ctx, a.ResultId)
	if err != nil {
		return "", err
	}
	if resultInfo.UserId != a.UserId {
		return "", errno.NewErrNo(errno.ServiceUserErrorAppealCode, "user have not permission to appeal the result")
	}
	// 申诉
	appeal_id, err := mysql.CreateAppeal(svc.ctx, a)
	if err != nil {
		return "", err
	}
	// 异步同步result内的信息
	taskqueue.AddAppealToScoreTask(svc.ctx, constants.AppealKey, a.ResultId, appeal_id, constants.OnAppeal)
	return appeal_id, nil
}

func (svc *AppealService) DeleteAppeal(appeal_id string) error {
	exist, err := mysql.IsAppealExistByAppealId(svc.ctx, appeal_id)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceAppealNotExistCode, "appeal not exist")
	}
	stu_id := GetUserIDFromContext(svc.c)
	// 检验appeal属于user
	appeal, err := mysql.QueryAppealById(svc.ctx, appeal_id)
	if err != nil {
		return err
	}
	if appeal.UserId != stu_id {
		return errno.NewErrNo(errno.ServiceUserErrorAppealCode, "user have not permission to delete appeal")
	}
	if appeal.Status != constants.WaitResult {
		return errno.NewErrNo(errno.ServiceAppealUnchangedCode, "appeal cannot be deleted after handle")
	}
	// 删除
	err = mysql.DeleteAppealById(svc.ctx, appeal_id)
	if err != nil {
		return err
	}
	// 异步进行清除
	taskqueue.AddAppealToScoreTask(svc.ctx, constants.AppealKey, appeal.ResultId, "0", constants.OffAppeal)
	return nil
}
func (svc *AppealService) QueryAppealById(appeal_id string) (*model.Appeal, error) {
	exist, err := mysql.IsAppealExistByAppealId(svc.ctx, appeal_id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceAppealNotExistCode, "appeal not exist")
	}
	stu_id := GetUserIDFromContext(svc.c)
	// 检验appeal属于user
	appeal, err := mysql.QueryAppealById(svc.ctx, appeal_id)
	if err != nil {
		return nil, err
	}
	if appeal.UserId != stu_id {
		return nil, errno.NewErrNo(errno.ServiceUserErrorAppealCode, "user have not permission to query appeal")
	}
	return appeal, nil
}

func (svc *AppealService) QueryStuAllAppeals() ([]*model.Appeal, int64, error) {
	stu_id := GetUserIDFromContext(svc.c)
	appeals, count, err := mysql.QueryAppealByUserId(svc.ctx, stu_id)
	if err != nil {
		return nil, -1, err
	}
	return appeals, count, nil
}

// 处理申诉
func (svc *AppealService) HandleAppeal(appeal *model.Appeal) error {
	// 常规检验存在性
	exist, err := mysql.IsAppealExistByAppealId(svc.ctx, appeal.AppealId)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceAppealNotExistCode, "appeal not exist")
	}
	id := GetUserIDFromContext(svc.c)
	appeal.UserId = id
	// 判读有无权限处理该申诉
	appealInfo, err := mysql.QueryAppealById(svc.ctx, appeal.AppealId)
	if err != nil {
		return err
	}
	if appealInfo.Status == "approved" || appealInfo.Status == "rejected" {
		return errno.NewErrNo(errno.ServiceRepeatAction, "the appeal have been handled")
	}
	exist, err = mysql.IsAdminRelationExist(svc.ctx, id, appealInfo.UserId)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceNoAuthToDo, "No permission to handle the stu's appeal")
	}
	err = mysql.UpdateAppealInfo(svc.ctx, appeal)
	if err != nil {
		return err
	}
	// 由于提供了修改接口，这边不在异步进行
	return nil
}
func (svc *AppealService) QueryBelongStuAppeal(status string) ([]*model.Appeal, int64, error) {
	if status != "pending" && status != "approved" && status != "rejected" {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "error status type")
	}
	user_id := GetUserIDFromContext(svc.c)
	stuList, err := mysql.QueryStuByAdmin(svc.ctx, user_id)
	if err != nil {
		return nil, -1, err
	}
	if len(stuList) == 0 {
		return nil, 0, nil
	}
	appealInfoList := make([]*model.Appeal, 0)
	var totalCount int64 = 0
	for _, v := range stuList {
		appeals, _, err := mysql.QueryAppealByUserId(svc.ctx, v) // events 是 []*model.Event
		if err != nil {
			return nil, -1, err
		}
		// 过滤符合状态的事件
		for _, appeal := range appeals {
			if appeal.Status == status {
				appealInfoList = append(appealInfoList, appeal)
				totalCount++
			}
		}
	}
	return appealInfoList, totalCount, nil
}
