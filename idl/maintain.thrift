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
// 添加用户
struct AddUserRequest{
    1: required string user_role,
    2: required string user_id,
    3: required string password,
    4: required string email,
    5: required string username,
    6: optional string college //可选添加学院信息
}
struct AddUserResponse{
    1: required model.BaseResp base,
    2: required string user_id,
}
struct AddAdminObjectRequest{
    1: required string user_id,
    2: optional string major_name,
    3: optional string grade,
    4: optional string college_name,
}
struct AddAdminObjectResponse{
    1: required model.BaseResp base,
}

struct UploadRewardFileRequest{
    1: binary data (api.form="reward_file")
}

struct UploadRewardFileResponse{
    1: required model.BaseResp base,
    2: required string event_id,
}

struct UploadRewardRequest {
    1: required string college,
    2: required string event_name,
    3: required string organizer,
    4: required string event_time,
    5: required string related_majors,
    6: required string applicable_majors,
    7: required string recognition_basis,
    8: required string recognized_level,
    9: required string is_active,
}

struct UploadRewardResponse {
    1: required model.BaseResp base,
    2: required string recognized_event_id,
}



service maintainService{
     QueryAllCollegeResponse QueryCollege(1: QueryAllCollegeRequest req) (api.get = "/api/admin/colleges"),
     QueryMajorByCollegeIdResponse QueryMajorByCollegeId(1: QueryMajorByCollegeIdRequest req) (api.get = "/api/admin/majors"),
     UploadMajorResponse UploadMajor(1: UploadMajorRequest req) (api.post = "/api/admin/majors"),
     UploadCollegeResponse UploadCollege(1: UploadCollegeRequest req) (api.post = "/api/admin/colleges"),
     AddUserResponse AddUser(1:AddUserRequest req)(api.post = "/api/admin/users"),
     AddAdminObjectResponse AddAdminObject(1:AddAdminObjectRequest req)(api.post = "/api/admin/users/permission"),
     UploadRewardResponse UploadReward(1:UploadRewardRequest req)(api.post = "/api/admin/upload/reward/single")
}

