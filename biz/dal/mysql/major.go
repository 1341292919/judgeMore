package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

type GetMajorInfoByCollegeIdFunc func(ctx context.Context, college_id int64) ([]*model.Major, int64, error)
type CreateMajorFunc func(ctx context.Context, major *model.Major) (int64, error)
type QueryMajorByIdFunc func(ctx context.Context, major_id int64) (*model.Major, error)
type IsMajorExistFunc func(ctx context.Context, majorId int64) (bool, error)
type DeleteMajorByIdFUnc func(ctx context.Context, majorId int64) error
type UpdateMajorFunc func(ctx context.Context, majorId int64, updateFields map[string]interface{}) (*model.Major, error)

// 对外暴露的函数变量（默认指向真实实现）
var (
	GetMajorInfoByCollegeId GetMajorInfoByCollegeIdFunc = RealGetMajorInfoByCollegeId
	CreateMajor             CreateMajorFunc             = RealCreateMajor
	QueryMajorById          QueryMajorByIdFunc          = RealQueryMajorById
	IsMajorExist            IsMajorExistFunc            = RealIsMajorExist
	DeleteMajorById         DeleteMajorByIdFUnc         = RealDeleteMajorById
	UpdateMajor             UpdateMajorFunc             = RealUpdateMajor
)

func RealGetMajorInfoByCollegeId(ctx context.Context, college_id int64) ([]*model.Major, int64, error) {
	var majorInfos []*Major
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("college_id = ?", college_id).
		Find(&majorInfos).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query stu event: %v", err)
	}
	return BuildMajorInfoList(majorInfos), count, err
}

func BuildMajorInfo(data *Major) *model.Major {
	return &model.Major{
		MajorId:   data.MajorId,
		MajorName: data.MajorName,
		CollegeId: data.CollegeId,
	}
}

func BuildMajorInfoList(data []*Major) []*model.Major {
	resp := make([]*model.Major, 0)
	for _, v := range data {
		s := BuildMajorInfo(v)
		resp = append(resp, s)
	}
	return resp
}

func RealCreateMajor(ctx context.Context, m *model.Major) (int64, error) {
	var major *Major
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("major_name = ?", m.MajorName).
		First(&major).
		Error
	if err == nil { //找到了
		return -1, errno.NewErrNo(errno.ServiceMajorExistCode, "major_name already exists")
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query major: %v", err)
		}
	}
	major.MajorName = m.MajorName
	major.CollegeId = m.CollegeId
	err = db.WithContext(ctx).
		Table(constants.TableMajor).
		Create(major).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create major: %v", err)
	}
	return major.MajorId, nil
}
func RealQueryMajorById(ctx context.Context, major_id int64) (*model.Major, error) {
	var majorInfo *Major
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("id = ?", major_id).
		Find(&majorInfo).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query major by majorid: %v", err)
	}
	return BuildMajorInfo(majorInfo), err
}
func RealIsMajorExist(ctx context.Context, majorId int64) (bool, error) {
	var majorInfo *Major
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("major_id = ?", majorId).
		First(&majorInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query major: %v", err)
	}
	return true, nil
}

func RealDeleteMajorById(ctx context.Context, majorId int64) error {
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("major_id = ?", majorId).
		Delete(&model.Major{}).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete major: %v", err)
	}
	return nil
}

func RealUpdateMajor(ctx context.Context, majorId int64, updateFields map[string]interface{}) (*model.Major, error) {
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("major_id = ?", majorId).
		Updates(updateFields).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update major: %v", err)
	}
	return QueryMajorById(ctx, majorId)
}
