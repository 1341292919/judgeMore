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

// 对外暴露的函数变量（默认指向真实实现）
var (
	GetMajorInfoByCollegeId GetMajorInfoByCollegeIdFunc = RealGetMajorInfoByCollegeId
	CreateMajor             CreateMajorFunc             = RealCreateMajor
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
func GetAllMajorInfo(ctx context.Context) ([]*model.Major, int64, error) {
	var majorInfos []*Major
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
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
