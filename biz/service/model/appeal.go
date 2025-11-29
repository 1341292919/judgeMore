package model

type Appeal struct {
	AppealId       string
	ResultId       string
	EventId        string
	EventName      string
	EventLevel     string
	AwardLevel     string
	Score          float64
	UserId         string
	AppealType     string
	AppealReason   string
	AttachmentPath string
	Status         string
	HandledBy      string
	HandledAt      int64
	HandleResult   string
	AppealCount    int64
	CreateAT       int64
	UpdateAT       int64
	DeleteAT       int64
}
