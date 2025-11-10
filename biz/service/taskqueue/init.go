package taskqueue

import (
	"judgeMore/pkg/constants"
	"judgeMore/pkg/taskqueue"
)

var taskQueue *taskqueue.BaseTaskQueue

func Init() {
	taskQueue = taskqueue.NewBaseTaskQueue()
	Work(constants.EventKey)
	Work(constants.AppealKey)
	Work(constants.StructKey)
}
