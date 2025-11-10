namespace go appeal
include "model.thrift"

struct ApplyAppealRequest{
    1 : required string result_id,
    2 : required string appeal_message
    3 : optional string attachment_path,
    4 : required string appeal_type
}
struct ApplyAppealResponse{
         1: required model.BaseResp base,
         2: required string appeal_id,
}
struct QueryAppealInfoRequest{
    1:required string appeal_id,
}
struct QueryAppealInfoResponse{
             1: required model.BaseResp base,
             2: required model.Appeal data,
}
struct QueryStuAppealInfoRequest{

}
struct QueryStuAppealInfoResponse{
             1: required model.BaseResp base,
             2: required model.AppealList data,
}
struct DeleteAppealRequest{
     1 : required string appeal_id,
}
struct DeleteAppealResponse{
      1: required model.BaseResp base,
}
struct UpdateAppealRequest{
      1 :required string appeal_id,
      2 :required string status,
      3 :required string handled_result,
}
struct UpdateAppealResponse{
      1: required model.BaseResp base,
}
struct QueryBelongStuAppealRequest{
    1: required string status,
}
struct QueryBelongStuAppealResponse{
     1: required model.BaseResp base,
     2: required model.AppealList data,
}
service AppealService{
    ApplyAppealResponse ApplyAppeal(1:ApplyAppealRequest req)(api.post = "/api/upload/appeal"),
    DeleteAppealResponse DeleteAppeal(1:DeleteAppealRequest req)(api.delete = "/api/delete/appeal"),
    QueryAppealInfoResponse QueryAppealInfo(1:QueryAppealInfoRequest req)(api.get="/api/query/appeal"),
    QueryStuAppealInfoResponse QueryStuAppealInfo(1:QueryStuAppealInfoRequest req)(api.get="/api/query/appeal/stu"),
    UpdateAppealResponse UpdateAppealStatus(1:UpdateAppealRequest req)(api.post = "/api/appeal/status"),
    QueryBelongStuAppealResponse QueryBelongStuAppeal(1:QueryBelongStuAppealRequest req)(api.get = "/api/admin/query/appeal/stu")
}