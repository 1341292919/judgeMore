package model

type College struct {
	CollegeId   int64
	CollegeName string
}

type Major struct {
	MajorId   int64
	MajorName string
	CollegeId int64
}
type Relation struct {
	RelationId  string
	UserId      string
	CollegeId   string
	CollegeName string
	MajorName   string
	MajorId     string
	Grade       string
}
