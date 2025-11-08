package model

type User struct {
	Uid      string
	UserName string
	Email    string
	Password string
	College  string
	Major    string
	Grade    string
	Status   int
	Role     string
	CreateAT int64
	UpdateAT int64
	DeleteAT int64
}

type EmailAuth struct {
	Code  string
	Email string
	Uid   string
	Time  int64 //时间戳
}
type UpdateUserRequest struct {
	UserId    int64
	UserName  string
	Email     string
	Password  string
	CollegeId int64
	MajorId   int64
	Grade     string
}

type UpdateMajorRequest struct {
	MajorId   int64
	CollegeId int64
	MajorName string
}

type UpdateCollegeRequest struct {
	CollegeId   int64
	CollegeName string
}
type QueryUserRequest struct {
	CollegeId int64
	MajorId   int64
	Role      string
}
