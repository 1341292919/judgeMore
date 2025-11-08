package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func IsUserExist(ctx context.Context, user *model.User) (bool, error) {
	var userInfo *User
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", user.Uid).
		First(&userInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了说明用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return true, nil
}
func CreateUser(ctx context.Context, user *model.User) (string, error) {
	userInfo := &User{
		UserName: user.UserName,
		Password: user.Password,
		Email:    user.Email,
		RoleId:   user.Uid,
		UserRole: "student",
		Status:   0, //初始状态未激活
	}
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Create(userInfo).
		Error
	if err != nil {
		return "", errno.NewErrNo(errno.InternalDatabaseErrorCode, "Create User Error:"+err.Error())
	}
	return userInfo.RoleId, nil
}

// 该函数调用前检验存在性
func GetUserInfoByRoleId(ctx context.Context, role_id string) (*model.User, error) {
	var userInfo *User
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", role_id).
		First(&userInfo).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query user Info error:"+err.Error())
	}
	return &model.User{
		Uid:      userInfo.RoleId,
		UserName: userInfo.UserName,
		Grade:    userInfo.Grade,
		Major:    userInfo.Major,
		College:  userInfo.College,
		Password: userInfo.Password,
		Status:   userInfo.Status,
		Email:    userInfo.Email,
		Role:     userInfo.UserRole,
		UpdateAT: userInfo.UpdatedAt.Unix(),
		CreateAT: userInfo.CreatedAt.Unix(),
		DeleteAT: 0,
	}, nil
}
func UpdateUserPassword(ctx context.Context, user *model.User) error {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", user.Uid).
		Update("password", user.Password).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "Update User Password"+err.Error())
	}
	return nil
}

func UpdateInfoByRoleId(ctx context.Context, role_id string, element ...string) (*model.User, error) {
	updateFields := make(map[string]interface{})
	for i, value := range element {
		if value == "" {
			continue // 跳过空值
		}
		switch i {
		case 0:
			updateFields["major"] = value
		case 1:
			updateFields["college"] = value
		case 2:
			updateFields["grade"] = value
		}
	}
	// 存在多值更新 以事务提交保证原子性
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table(constants.TableUser).
			Where("role_id = ?", role_id).
			Updates(updateFields).
			Error
	})
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update userInfo: %v", err)
	}

	return GetUserInfoByRoleId(ctx, role_id)
}

func ActivateUser(ctx context.Context, uid string) error {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", uid).
		Update("status", 1).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to activate user: %v", err)
	}
	return nil
}
// 更新用户状态
func UpdateUserStatus(ctx context.Context, uid string, status int64) (*model.User, error) {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", uid).
		Update("status", status).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update user status: %v", err)
	}
	return GetUserInfoByRoleId(ctx, uid)
}

func QueryUserByCondition(ctx context.Context, page_size, page_num int64, req *model.QueryUserRequest) ([]*model.User, int64, error) {

	tx := db.WithContext(ctx).Table(constants.TableUser)
	// 拼接查询条件（仅当参数有效时添加条件）
	if req.CollegeId > 0 {
		tx = tx.Where("college_id = ?", req.CollegeId)
	}
	if req.MajorId > 0 {
		tx = tx.Where("major_id = ?", req.MajorId)
	}
	if req.Role != "" {
		tx = tx.Where("role = ?", strings.TrimSpace(req.Role)) // 去除空格避免无效查询
	}
	// 先查询总条数（不带分页的条件计数）
	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to count user: %v", err)
	}
	// 处理分页（计算偏移量，添加LIMIT和OFFSET）
	// pageNum从1开始，偏移量 = (pageNum-1) * pageSize
	offset := (page_num - 1) * page_size
	tx = tx.Offset(int(offset)).Limit(int(page_size))
	// 执行分页查询
	var users []*User
	err2 := tx.Find(&users).Error
	if err2 != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return BuildUserInfoList(users), total, nil
}
func UpdateUser(ctx context.Context, uid string, updateFields map[string]interface{}) (*model.User, error) {

	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("uid = ?", uid).
		Updates(updateFields).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: update user failed: %v", err)
	}
	return GetUserInfoByRoleId(ctx, uid)
}
func UploadUser(ctx context.Context, uploadFields map[string]interface{}) (*model.User, error) {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Create(uploadFields).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: upload user failed: %v", err)
	}
	return GetUserInfoByRoleId(ctx, uploadFields["role_id"].(string))

}
func BuildUserInfo(data *User) *model.User {
	return &model.User{
		Uid:      data.RoleId,
		UserName: data.UserName,
		Grade:    data.Grade,
		Major:    data.Major,
		College:  data.College,
		Password: data.Password,
		Status:   data.Status,
		Email:    data.Email,
		Role:     data.UserRole,
		UpdateAT: data.UpdatedAt.Unix(),
		CreateAT: data.CreatedAt.Unix(),
		DeleteAT: 0,
	}
}

func BuildUserInfoList(data []*User) []*model.User {
	resp := make([]*model.User, 0)
	for _, v := range data {
		s := BuildUserInfo(v)
		resp = append(resp, s)
	}
	return resp
}
