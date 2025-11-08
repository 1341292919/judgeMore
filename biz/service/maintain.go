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
// 激活、禁用用户
func (svc *MaintainService) UpdateUserStatus(user_id int64, status int64) (*model.User, error) {
	//检查user是否存在
	uid := strconv.FormatInt(user_id, 10)
	userInfo := &model.User{
		Uid: uid,
	}
	exist, err := mysql.IsUserExist(svc.ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("check user exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "user not exist")
	}
	// 检验传上来的status
	if status != 1 && status != 0 {
		return nil, fmt.Errorf("status should be 0 or 1")
	}
	//更新user的状态
	info, err := mysql.UpdateUserStatus(svc.ctx, uid, status)
	if err != nil {
		return nil, fmt.Errorf("update euser status failed: %w", err)
	}
	return info, nil
}

// 检查UpdatreUserRequest
func (svc *MaintainService) ExamineUpdateUserRequest(req *model.UpdateUserRequest) (map[string]interface{}, error) {
	//检查userid是否存在
	if req.UserId <= 0 {
		return nil, errno.NewErrNo(errno.ParamMissingErrorCode, "user_id is required")
	}
	uid := strconv.FormatInt(req.UserId, 10)
	//检查该学生是否存在
	userInfo := &model.User{
		Uid: uid,
	}
	exist, err := mysql.IsUserExist(svc.ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("check user exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "user not exist")
	}
	updateFields := make(map[string]interface{})
	if req.UserName != "" {
		updateFields["user_name"] = req.UserName
	}
	if req.CollegeId > 0 {
		collegeInfo, err2 := mysql.QueryCollegeById(svc.ctx, req.CollegeId)
		if err2 != nil {
			return nil, fmt.Errorf("query College failed: %w", err)
		}
		updateFields["college"] = collegeInfo.CollegeName
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	if req.Password != "" {
		updateFields["password"] = req.Password
	}
	if req.Grade != "" {
		updateFields["grade"] = req.Grade
	}
	if req.MajorId > 0 {
		majorInfo, err3 := mysql.QueryMajorById(svc.ctx, req.MajorId)
		if err3 != nil {
			return nil, fmt.Errorf("query major failed: %w", err3)
		}
		updateFields["major"] = majorInfo.MajorName
	}
	updateFields["updated_at"] = time.Now()
	return updateFields, nil
}

// 更新用户信息
func (svc *MaintainService) UpdateUser(req *model.UpdateUserRequest) (*model.User, error) {

	//将需要更新的信息包装成map
	updateFields, err := svc.ExamineUpdateUserRequest(req)
	if err != nil {
		return nil, errno.NewErrNo(errno.ParamLogicalErrorCode, "UpdateUserRequest is not standardized")
	}
	uid := strconv.FormatInt(req.UserId, 10)
	updatedUser, err2 := mysql.UpdateUser(svc.ctx, uid, updateFields)
	if err2 != nil {
		return nil, fmt.Errorf("update user failed: %w", err)
	}
	// 转换数据库模型为返回的 UserInfo
	userInfo2 := &model.User{
		Uid:      uid,
		UserName: updatedUser.UserName,
		College:  updatedUser.College,
		Major:    updatedUser.Major,
		Grade:    updatedUser.Grade,
		Email:    updatedUser.Email,
		Password: updatedUser.Password,
		CreateAT: updatedUser.CreateAT,
		UpdateAT: updatedUser.UpdateAT,
	}

	return userInfo2, nil
}

// 上传新的用户信息
func (svc *MaintainService) UploadUser(req *model.UpdateUserRequest) (*model.User, error) {

	//将学生信息包装成map，但此时缺少role_id
	uploadFields, err := svc.ExamineUpdateUserRequest(req)
	if err != nil {
		return nil, errno.NewErrNo(errno.ParamLogicalErrorCode, "UpdateUserRequest is not standardized")
	}
	uploadUser, err2 := mysql.UploadUser(svc.ctx, uploadFields)
	if err2 != nil {
		return nil, fmt.Errorf("update user failed: %w", err)
	}
	//添加role_id
	uid := strconv.FormatInt(req.UserId, 10)
	uploadFields["role_id"] = uid
	//管理员创建的用户，默认已激活
	uploadFields["status"] = 1
	//存入数据库
	uploadUser, err3 := mysql.UploadUser(svc.ctx, uploadFields)
	if err3 != nil {
		return nil, fmt.Errorf("upload user failed: %w", err3)
	}
	// 转换数据库模型为返回的 UserInfo
	userInfo2 := &model.User{
		Uid:      uid,
		UserName: uploadUser.UserName,
		College:  uploadUser.College,
		Major:    uploadUser.Major,
		Grade:    uploadUser.Grade,
		Email:    uploadUser.Email,
		Password: uploadUser.Password,
		Status:   uploadUser.Status,
		CreateAT: uploadUser.CreateAT,
		UpdateAT: uploadUser.UpdateAT,
	}

	return userInfo2, nil
}

// 查询所有用户信息
func (svc *MaintainService) QueryUserInfo(page_num, page_size int64, req *model.QueryUserRequest) ([]*model.User, int64, error) {
	//  参数校验（避免无效查询）
	if page_num <= 0 {
		return nil, 0, fmt.Errorf("invalid page_num: %d", page_num)
	}
	if page_size <= 0 || page_size > 1000 { // 限制最大页大小，防止恶意查询
		return nil, 0, fmt.Errorf("invalid page_size: %d", page_size)
	}

	// 2调用数据库层带条件的分页查询
	users, total, err := mysql.QueryUserByCondition(svc.ctx, page_num, page_size, req)
	if err != nil {
		return nil, 0, fmt.Errorf("query user failed: %w", err)
	}
	// 3. 返回结果（数据库已分页，直接返回）
	return users, total, nil
}

// 删除专业
func (svc *MaintainService) DeleteMajor(majorId int64) error {
	//检查
	exist, err := mysql.IsMajorExist(svc.ctx, majorId)
	if err != nil {
		return fmt.Errorf("check major exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceEventExistCode, "major not exist")
	}
	//专业存在，删除专业
	err2 := mysql.DeleteMajorById(svc.ctx, majorId)
	if err2 != nil {
		return fmt.Errorf("delete major failed: %w", err)
	}
	return nil
}

// 删除学院 只有当专业为空时，学院才能删除
func (svc *MaintainService) DeleteCollege(collegeId int64) error {
	//检查
	exist, err := mysql.IsMajorExist(svc.ctx, collegeId)
	if err != nil {
		return fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceEventExistCode, "college not exist")
	}
	//检查此学院的专业是否为空
	_, count, err2 := mysql.GetMajorInfoByCollegeId(svc.ctx, collegeId)
	if err2 != nil {
		return fmt.Errorf("check major exist failed: %w", err)
	}
	if count != 0 {
		return errno.NewErrNo(errno.ServiceEventExistCode, "major exist")
	}
	//不为空，删除学院
	err3 := mysql.DeleteCollegeById(svc.ctx, collegeId)
	if err3 != nil {
		return fmt.Errorf("delete college failed: %w", err)
	}
	return nil
}

// 更新专业
func (svc *MaintainService) UpdateMajor(req model.UpdateMajorRequest) (*model.Major, error) {
	//检查专业是否存在
	exist, err := mysql.IsMajorExist(svc.ctx, req.MajorId)
	if err != nil {
		return nil, fmt.Errorf("check major exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "major not exist")
	}
	//组合更新信息为map
	updateFields := make(map[string]interface{})
	if req.MajorName != "" {
		updateFields["user_name"] = req.MajorName
	}
	if req.CollegeId != 0 {
		updateFields["college_id"] = req.CollegeId
	}
	//更新数据库
	majorInfo, err2 := mysql.UpdateMajor(svc.ctx, req.MajorId, updateFields)
	if err2 != nil {
		return nil, fmt.Errorf("update major failed: %w", err2)
	}
	return majorInfo, nil
}

// 更新学院
func (svc *MaintainService) UpdateCollege(req model.UpdateCollegeRequest) (*model.College, error) {
	//检查
	exist, err := mysql.IsCollegeExist(svc.ctx, req.CollegeId)
	if err != nil {
		return nil, fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "college not exist")
	}
	//组合更新信息为map
	updateFields := make(map[string]interface{})
	if req.CollegeName != "" {
		updateFields["user_name"] = req.CollegeName
	}
	//更新数据库
	collegeInfo, err2 := mysql.UpdateCollege(svc.ctx, req.CollegeId, updateFields)
	if err2 != nil {
		return nil, fmt.Errorf("update college failed: %w", err2)
	}
	return collegeInfo, nil
}
