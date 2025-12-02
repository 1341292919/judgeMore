package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"sort"
	"strconv"
)

func ScoreRecord(data *model.ScoreRecord) *resp.ScoreRecord {
	return &resp.ScoreRecord{
		ScoreID:    data.ResultId,
		EventID:    data.EventId,
		UserID:     data.UserId,
		RuleID:     data.RuleId,
		AppealID:   data.AppealId,
		FinalScore: data.FinalIntegral,
		Status:     data.Status,
		CreatedAt:  strconv.FormatInt(data.CreateAT, 10),
		UpdatedAt:  strconv.FormatInt(data.UpdateAT, 10),
		DeletedAt:  strconv.FormatInt(data.DeleteAT, 10),
	}
}

func ScoreRecordList(data []*model.ScoreRecord, total int64) *resp.ScoreRecordList {
	info := make([]*resp.ScoreRecord, 0)
	var sum float64
	for _, v := range data {
		sum += v.FinalIntegral
		info = append(info, ScoreRecord(v))
	}
	return &resp.ScoreRecordList{
		Items: info,
		Total: total,
		Sum:   sum,
	}
}

func StuScoreMessage(data *model.StuScoreMessage) *resp.StuScoreMessage {
	return &resp.StuScoreMessage{
		StuID:   data.Uid,
		StuName: data.Name,
		College: data.College,
		Grade:   data.Grade,
		Score:   data.Score,
	}
}
func StuScoreMessageList(data []*model.StuScoreMessage, c int64) *resp.StuScoreMessageList {
	info := make([]*resp.StuScoreMessage, 0, len(data))
	for _, v := range data {
		info = append(info, StuScoreMessage(v))
	}

	// 按 Score 从高到低排序
	sort.Slice(info, func(i, j int) bool {
		return info[i].Score > info[j].Score
	})

	return &resp.StuScoreMessageList{
		Item:  info,
		Total: c,
	}
}
