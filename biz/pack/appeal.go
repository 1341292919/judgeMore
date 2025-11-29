package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"strconv"
)

func Appeal(data *model.Appeal) *resp.Appeal {
	return &resp.Appeal{
		UserID:         data.UserId,
		AppealID:       data.AppealId,
		ResultID:       data.ResultId,
		EventID:        data.EventId,
		EventName:      data.EventName,
		AwardLevel:     data.AwardLevel,
		FinalScore:     data.Score,
		HandleResult:   data.HandleResult,
		HandleBy:       data.HandledBy,
		AppealReason:   data.AppealReason,
		AppealCount:    data.AppealCount,
		AttachmentPath: data.AttachmentPath,
		Status:         data.Status,
		AppealType:     data.AppealType,
		HandleTime:     strconv.FormatInt(data.HandledAt, 10),
		CreatedAt:      strconv.FormatInt(data.CreateAT, 10),
		UpdatedAt:      strconv.FormatInt(data.UpdateAT, 10),
		DeletedAt:      strconv.FormatInt(data.DeleteAT, 10),
	}
}

func AppealList(data []*model.Appeal, total int64) *resp.AppealList {
	res := make([]*resp.Appeal, 0)
	for _, v := range data {
		res = append(res, Appeal(v))
	}
	return &resp.AppealList{
		Items: res,
		Total: total,
	}
}
