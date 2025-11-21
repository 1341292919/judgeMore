package service

import (
	"context"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/es"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
)

// 查询所有认可奖项的函数
func QueryAllRecognizedReward(ctx context.Context) ([]*model.RecognizedEvent, error) {
	exist, err := cache.IsRecognizeEventExist(ctx)
	if err != nil {
		return nil, err
	}
	var recognizedEventList []*model.RecognizedEvent
	if !exist {
		// db 载入 redis
		recognizedEventList, _, err = mysql.QueryRecognizedEvent(ctx)
		if err != nil {
			return nil, err
		}
		err = cache.RecognizeEventToCache(ctx, recognizedEventList)
		if err != nil {
			return nil, err
		}
	} else {
		recognizedEventList, err = cache.QueryAllRecognizeEvent(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 筛去isactive == false的
	filteredList := make([]*model.RecognizedEvent, 0, len(recognizedEventList))
	for _, v := range recognizedEventList {
		if v.IsActive {
			filteredList = append(filteredList, v)
		}
	}
	return filteredList, nil
}

func IsRecognizedEventExist(ctx context.Context, recognizde_id string) (bool, error) {
	re, err := QueryAllRecognizedReward(ctx)
	if err != nil {
		return false, nil
	}
	for _, m := range re {
		if m.RecognizedEventId == recognizde_id && m.IsActive == true {
			return true, nil
		}
	}
	return false, err
}

// 根据search_req来从es搜索奖项
func SearchRecognizedEvent(ctx context.Context, req *model.ViewRecognizedRewardReq) ([]*model.RecognizedEvent, error) {
	exist, err := es.IsIndexDataExist(ctx, constants.IndexName)
	if err != nil {
		return nil, err
	}
	var reList []*model.RecognizedEvent
	if !exist {
		reList, _, err = mysql.QueryRecognizedEvent(ctx)
		for _, v := range reList {
			err = es.AddItem(ctx, constants.IndexName, v)
			if err != nil {
				return nil, err
			}
		}
	}
	result, _, err := es.SearchItems(ctx, constants.IndexName, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}
