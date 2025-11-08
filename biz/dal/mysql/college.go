package mysql

import (
	"context"
	"errors"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"

	"gorm.io/gorm"
)

type GetCollegeInfoFunc func(ctx context.Context) ([]*model.College, int64, error)
type IsCollegeExistFunc func(ctx context.Context, college_id int64) (bool, error)
type CreateNewCollegeFunc func(ctx context.Context, college_name string) (int64, error)
type QueryCollegeByIdFunc func(ctx context.Context, college_id int64) (*model.College, error)
type DeleteCollegeByIdFunc func(ctx context.Context, collegeId int64) error
type UpdateCollegeFunc func(ctx context.Context, CollegeId int64, updateFields map[string]interface{}) (*model.College, error)


// 对外暴露的函数变量（默认指向真实实现,用于测试）
var (
	GetCollegeInfo   GetCollegeInfoFunc   = RealGetCollegeInfo
	IsCollegeExist   IsCollegeExistFunc   = RealIsCollegeExist
	CreateNewCollege CreateNewCollegeFunc = RealCreateNewCollege
	QueryCollegeById  QueryCollegeByIdFunc  = RealQueryCollegeById
	DeleteCollegeById DeleteCollegeByIdFunc = RealDeleteCollegeById
	UpdateCollege     UpdateCollegeFunc     = RealUpdateCollege
)

func RealGetCollegeInfo(ctx context.Context) ([]*model.College, int64, error) {
	var collegeInfos []*College
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Find(&collegeInfos).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query college: %v", err)
	}
	return BuildCollegeInfoList(collegeInfos), count, err
}

func BuildCollegeInfo(data *College) *model.College {
	return &model.College{
		CollegeId:   data.CollegeId,
		CollegeName: data.CollegeName,
	}
}
func BuildCollegeInfoList(data []*College) []*model.College {
	resp := make([]*model.College, 0)
	for _, v := range data {
		s := BuildCollegeInfo(v)
		resp = append(resp, s)
	}
	return resp
}
func RealIsCollegeExist(ctx context.Context, college_id int64) (bool, error) {
	var collegeInfo *College
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_id = ?", college_id).
		First(&collegeInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query college: %v", err)
	}
	return true, nil
}
func RealCreateNewCollege(ctx context.Context, college_name string) (int64, error) {
	var college *College
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_name = ?", college_name).
		First(&college).
		Error
	if err == nil { //找到了
		return -1, errno.NewErrNo(errno.ServiceCollegeExistCode, "college have exist")
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query college: %v", err)
		}
	}
	c := &College{
		CollegeName: college_name,
	}
	err = db.WithContext(ctx).
		Table(constants.TableCollege).
		Create(&c).
		Error
	if err != nil {
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create college: %v", err)
	}
	return c.CollegeId, nil
}
func RealQueryCollegeById(ctx context.Context, college_id int64) (*model.College, error) {
	var collegeInfo *College

	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_id = ?", college_id).
		First(&collegeInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return nil, nil
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query college: %v", err)
	}
	return BuildCollegeInfo(collegeInfo), nil
}

func RealDeleteCollegeById(ctx context.Context, collegeId int64) error {
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_id = ?", collegeId).
		Delete(&College{}).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete college: %v", err)
	}
	return nil
}

func RealUpdateCollege(ctx context.Context, CollegeId int64, updateFields map[string]interface{}) (*model.College, error) {
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_id = ?", CollegeId).
		Updates(updateFields).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update college: %v", err)
	}
	return QueryCollegeById(ctx, CollegeId)
}
