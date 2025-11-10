package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func CreateNewRelation(ctx context.Context, relation *model.Relation) error {
	var r = &Relation{
		UserId: relation.UserId,
	}
	if relation.CollegeId != "" {
		r.College = relation.CollegeName
		r.CollegeId = relation.CollegeId
	} else {
		r.MajorId = relation.MajorId
		r.Major = relation.MajorName
		r.Grade = relation.Grade
	}
	err := db.WithContext(ctx).
		Table(constants.TableRelation).
		Create(&r).
		Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "Create relation error"+err.Error())
	}
	return nil
}
func QueryAllRelation(ctx context.Context) ([]*model.Relation, int64, error) {
	var r []*Relation
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableRelation).
		Find(&r).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query relation error")
	}
	return buildRelationList(r), count, nil
}
func QueryRelationByUserId(ctx context.Context, user_id string) ([]*model.Relation, int64, error) {
	var r []*Relation
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableRelation).
		Where("user_id = ?", user_id).
		Find(&r).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query relation error")
	}
	return buildRelationList(r), count, nil
}
func relation(data *Relation) *model.Relation {
	return &model.Relation{
		RelationId:  data.RelationId,
		CollegeId:   data.CollegeId,
		MajorId:     data.MajorId,
		UserId:      data.UserId,
		CollegeName: data.College,
		MajorName:   data.Major,
		Grade:       data.Grade,
	}
}
func InsertAdminStu(ctx context.Context, user_id string, stuIdList []string) error {
	r := &AdminRelation{
		AdminId: user_id,
	}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, stu := range stuIdList {
			r.StuId = stu
			err := tx.
				Table(constants.TableAdminRelation).
				Create(&r).
				Error
			if err != nil {
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "create Admin Relation error:"+err.Error())
			}
		}
		return nil
	})
	return err
}
func InsertStuIdToAdminRelation(ctx context.Context, stu_id string, adminIdList []string) error {
	r := &AdminRelation{
		StuId: stu_id,
	}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, admin := range adminIdList {
			r.AdminId = admin
			err := tx.
				Table(constants.TableAdminRelation).
				Create(&r).
				Error
			if err != nil {
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "create Admin Relation error:"+err.Error())
			}
		}
		return nil
	})
	return err
}
func QueryStuByAdmin(ctx context.Context, admin_id string) ([]string, error) {
	var stuIds []string
	err := db.WithContext(ctx).
		Table(constants.TableAdminRelation).
		Where("admin_id = ?", admin_id).
		Pluck("stu_id", &stuIds).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query admin’s stu_id error"+err.Error())
	}

	return stuIds, nil
}
func IsAdminRelationExist(ctx context.Context, admin_id, stu_id string) (bool, error) {
	var Info *AdminRelation
	err := db.WithContext(ctx).
		Table(constants.TableAdminRelation).
		Where("admin_id = ? and stu_id = ?", admin_id, stu_id).
		First(&Info).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query admin relation: %v", err)
	}
	return true, nil
}
func buildRelationList(data []*Relation) []*model.Relation {
	result := make([]*model.Relation, 0)
	for _, r := range data {
		result = append(result, relation(r))
	}
	return result
}
