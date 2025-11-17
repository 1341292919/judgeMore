package mysql

import (
	"context"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
	"time"
)

func QueryRecognizedEvent(ctx context.Context) ([]*model.RecognizedEvent, int64, error) {
	var reconize_event []*RecognizedEvent
	var total int64
	err := db.WithContext(ctx).
		Table(constants.TableReconizedEvent).
		Find(&reconize_event).
		Count(&total).
		Error
	if err != nil {
		return nil, -1, errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to query reconized_event"+err.Error())
	}
	return buildRecognizedList(reconize_event), total, nil
}
func AddRecognizedEvent(ctx context.Context, re *model.RecognizedEvent) (*model.RecognizedEvent, error) {
	isactive := 1
	if re.IsActive == false {
		isactive = 0
	}
	recognizedEvent := &RecognizedEvent{
		RecognizedLevel:     re.RecognizedLevel,
		RecognizedEventName: re.RecognizedEventName,
		RecognizedEventTime: re.RecognizedEventTime,
		RecognitionBasis:    re.RecognitionBasis,
		RelatedMajors:       re.RelatedMajors,
		ApplicableMajors:    re.ApplicableMajors,
		IsActive:            isactive,
		Organizer:           re.Organizer,
		College:             re.College,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err := db.WithContext(ctx).
		Table(constants.TableReconizedEvent).
		Create(&recognizedEvent).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to create new recgnized event"+err.Error())
	}
	return buildRecognizeEvent(recognizedEvent), nil
}
func DeleteRecognized(ctx context.Context, id string) error {
	err := db.WithContext(ctx).
		Table(constants.TableReconizedEvent).
		Where("recognized_event_id = ?", id).
		Update("is_active", 0).
		Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to delete recgnized event:"+err.Error())
	}
	return nil
}

func buildRecognizeEvent(data *RecognizedEvent) *model.RecognizedEvent {
	isactive := true
	if data.IsActive == 0 {
		isactive = false
	}
	d := &model.RecognizedEvent{
		RecognizedEventId:   data.RecognizedEventId,
		RecognizedLevel:     data.RecognizedLevel,
		RecognizedEventName: data.RecognizedEventName,
		RecognizedEventTime: data.RecognizedEventTime,
		RecognitionBasis:    data.RecognitionBasis,
		College:             data.College,
		Organizer:           data.Organizer,
		RelatedMajors:       data.RelatedMajors,
		ApplicableMajors:    data.ApplicableMajors,
		IsActive:            isactive,
		UpdateAT:            data.UpdatedAt.Unix(),
		CreateAT:            data.CreatedAt.Unix(),
		DeleteAT:            0,
	}
	return d
}
func buildRecognizedList(data []*RecognizedEvent) []*model.RecognizedEvent {
	r := make([]*model.RecognizedEvent, 0)
	for _, v := range data {
		r = append(r, buildRecognizeEvent(v))
	}
	return r
}
