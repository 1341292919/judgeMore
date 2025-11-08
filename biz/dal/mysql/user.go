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
