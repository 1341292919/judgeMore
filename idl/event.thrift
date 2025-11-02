namespace go event
include "./model.thrift"

struct UploadEventFileRequest{
     1: binary data (api.form="file"),
}
struct UploadEventFileResponse{
     1: required model.BaseResp base,
     2: required string event_id,
}

struct UploadEventRequest{
     1: required string race_name,
     2: required string race_sponsor,
     3: required string race_time,
}
struct UploadEventResponse{
     1: required model.BaseResp base,
     2: required string event_id,
}

struct QueryEventRequest{
    1: required string event_id,
}
struct QueryEventResponse{
      1: required model.BaseResp base,
      2: required model.Event data,
}
struct QueryEventByStuIdRequest{
    1: required string page_num,
    2: required string page_size,
}
struct QueryEventByStuIdResponse{
     1: required model.BaseResp base,
     2: required model.EventList data,
}
struct UpdateEventStatusRequest{
    1: required string event_id,
    2: required i64 examine_status,
}
struct UpdateEventStatusResponse{
     1: required model.BaseResp base,
     2: required model.Event data,
}
service EventService{
    UploadEventFileResponse UploadEventFile(1: UploadEventFileRequest req) (api.post = "/api/update/materials/file"),
    UploadEventResponse UploadEvent(1: UploadEventRequest req) (api.post = "/api/update/materials/auto"),
    QueryEventResponse QueryEvent(1: QueryEventRequest req) (api.get = "/api/query/materials/single"),
    QueryEventByStuIdResponse QueryEventByStuIdFile(1: QueryEventByStuIdRequest req) (api.get = "/api/query/materials/stu"),
    UpdateEventStatusResponse UpdateEventStatus(1: UpdateEventStatusRequest req) (api.post = "/api/examine/materials"),
}