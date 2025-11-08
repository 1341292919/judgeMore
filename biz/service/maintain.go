package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/utils"
)

type MaintainService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewMaintainService(ctx context.Context, c *app.RequestContext) *MaintainService {
	return &MaintainService{
		ctx: ctx,
		c:   c,
	}
}

// 查找所有学院的信息
func (svc *MaintainService) QueryColleges(page_num, page_size int64) ([]*model.College, int64, error) {
	if utils.VerifyPageParam(page_num, page_size) {
		return nil, -1, errno.NewErrNo(errno.ParamVerifyErrorCode, "Page Param invalid")
	}
	collegeInfoList, count, err := mysql.GetCollegeInfo(svc.ctx)
	if err != nil {
		return nil, count, err
	}
	// 分页返回
	count = int64(len(collegeInfoList))
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if startIndex > count {
		return nil, 0, nil
	}
	if endIndex > count {
		endIndex = count
	}
	return collegeInfoList[startIndex:endIndex], count, nil
}

func (svc *MaintainService) QueryMajorByCollegeId(college_id int64, page_num, page_size int64) ([]*model.Major, int64, error) {
	if utils.VerifyPageParam(page_num, page_size) {
		return nil, -1, errno.NewErrNo(errno.ParamVerifyErrorCode, "Page Param invalid")
	}
	// 存在性检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return nil, -1, fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return nil, -1, errno.NewErrNo(errno.ServiceCollegeNotExistCode, "college not exist")
	}
	majorInfoList, count, err := mysql.GetMajorInfoByCollegeId(svc.ctx, college_id)
	if err != nil {
		return nil, count, err
	}
	// 分页返回
	count = int64(len(majorInfoList))
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if startIndex > count {
		return nil, 0, nil
	}
	if endIndex > count {
		endIndex = count
	}
	return majorInfoList[startIndex:endIndex], count, nil
}

func (svc *MaintainService) UploadMajor(major_name string, college_id int64) (int64, error) {
	// 检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return -1, err
	}
	if !exist {
		return -1, errno.NewErrNo(errno.ServiceEventNotExistCode, "college not exist")
	}
	// 保存到数据库
	major_id, err := mysql.CreateMajor(svc.ctx, &model.Major{MajorName: major_name, CollegeId: college_id})
	if err != nil {
		return -1, err
	}
	// 返回数据库生成的自增ID
	return major_id, nil
}
func (svc *MaintainService) UploadCollege(collegeName string) (int64, error) {
	collegeId, err := mysql.CreateNewCollege(svc.ctx, collegeName)
	if err != nil {
		return -1, err
	}
	// 返回数据库生成的自增ID
	return collegeId, nil
}
