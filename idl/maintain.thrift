namespace go maintain
include "./model.thrift"

// 这个用于获取所有学院信息
struct QueryAllCollegeRequest{
      1: required i64 page_num,
      2: required i64 page_size,
}
struct QueryAllCollegeResponse{
     1: required model.BaseResp base,
     2: required model.CollegeList data,
}
// 获取学院的专业
struct QueryMajorByCollegeIdRequest{
    1: required i64 page_num,
    2: required i64 page_size,
    3: required i64 college_id,
}
struct QueryMajorByCollegeIdResponse{
     1: required model.BaseResp base,
     2: required model.MajorList data,
}
// 上传专业
struct UploadMajorRequest{
     1: required string major_name,
     2: required i64 college_id,
}
struct UploadMajorResponse{
     1: required model.BaseResp base,
     2: required i64 major_id,
}
// 上传学院
struct UploadCollegeRequest{
     1: required string college_name,
}
struct UploadCollegeResponse{
     1: required model.BaseResp base,
     2: required i64 college_id,
}
// 修改学院信息
struct UpdateCollegeRequest{
    1: optional string college_name,
}
struct UpdateCollegeResponse{
    1: required model.BaseResp base,
    2: required model.College data,
}
// 修改专业信息
struct UpdateMajorRequest{
    1: optional string major_name,
    2: optional i64 college_id,
}
struct UpdateMajorResponse{
    1: required model.BaseResp base,
    2: required model.Major data,
}
//删除学院 - 只有在学院下没有关联专业的情况下才能删除学院
struct DeleteCollegeRequest{
    1: required i64 college_id,
}
struct DeleteCollegeResponse{
    1: required model.BaseResp base,
}
//删除专业
struct DeleteMajorRequest{
    1: required i64 major_id,
}
struct DeleteMajorResponse{
    1: required model.BaseResp base,
}

// 获取用户信息
struct QueryUserInfoRequest{
    1: required i64 page_num,
    2: required i64 page_size,
    3: optional i64 college_id,
    4: optional i64 major_id,
    5: optional string role,
}
struct QueryUserInfoResponse{
    1: required model.BaseResp base,
    2: required model.UserInfoList data,
}
// 上传用户
struct UploadUserRequest{
    1: required model.UserInfo data,
}
struct UploadUserResponse{
    1: required model.BaseResp base,
    2: required i64 user_id,
}
// 修改用户信息
struct UpdateUserRequest{
    1: required i64 user_id,
    2: optional string user_name,
    3: optional i64 college_id,
    4: optional i64 major_id,
    5: optional string grade,
    6: optional string email,
    7: optional string password,
}
struct UpdateUserResponse{
    1: required model.BaseResp base,
    2: required model.UserInfo data,
}
// 激活、禁用用户
struct UpdateUserStatusRequest{
    1: required i64 user_id,
    2: required i64 status,
}
struct UpdateUserStatusResponse{
    1: required model.BaseResp base,
    2: required model.UserInfo data,
}
//删除用户
struct DeleteUserRequest{
    1: required i64 user_id,
}
struct DeleteUserResponse{
    1: required model.BaseResp base,
}

service maintainService{
     QueryAllCollegeResponse QueryCollege(1: QueryAllCollegeRequest req) (api.get = "/api/admin/colleges"),
     QueryMajorByCollegeIdResponse QueryMajorByCollegeId(1: QueryMajorByCollegeIdRequest req) (api.get = "/api/admin/majors"),
     UploadMajorResponse UploadMajor(1: UploadMajorRequest req) (api.post = "/api/admin/majors"),
     UploadCollegeResponse UploadCollege(1: UploadCollegeRequest req) (api.post = "/api/admin/colleges"),
     UpdateUserStatusResponse UpdateUserStatus(1: UpdateUserStatusRequest req) (api.post = "/api/admin/user/status"),
     UpdateUserResponse UpdateUser (1: UpdateUserRequest req) (api.put="/api/admin/user"),
     UploadUserResponse UploadUser (1: UploadUserRequest req) (api.post="/api/admin/user"),
     QueryUserInfoResponse QueryUserInfo (1: QueryUserInfoRequest req) (api.get="/api/admin/user"),
     DeleteMajorResponse DeleteMajor (1: DeleteMajorRequest req) (api.delete="/api/admin/delete/major"),
     DeleteCollegeResponse DeleteCollege (1: DeleteCollegeRequest req) (api.delete="/api/admin/delete/college"),
     UpdateMajorResponse UpdateMajor (1: UpdateMajorRequest req) (api.put="/api/admin/majors"),
     UpdateCollegeResponse UpdateCollege (1: UpdateCollegeRequest req) (api.put="/api/admin/colleges"),
}

