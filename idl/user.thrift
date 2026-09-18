namespace go user

include "model.thrift"

struct GetUserInfoReq{
}

struct GetUserInfoResp{
    1: model.BaseResp base,
    2: model.UserInfo user_info,
}

struct UploadAvatarReq{
    1: string avatar,
}

struct UploadAvatarResp{
    1: model.BaseResp base,
}

struct UpdateSignatureReq{
    1: string signature,
}

struct UpdateSignatureResp{
    1: model.BaseResp base,
}

struct CreateUserReq{
    1: string user_id,
    2: string username,
}

struct CreateUserResp{
    1: model.BaseResp base,
}

service UserService{
    GetUserInfoResp GetUserInfo(1: GetUserInfoReq req),
    UploadAvatarResp UploadAvatar(1: UploadAvatarReq req),
    UpdateSignatureResp UpdateSignature(1: UpdateSignatureReq req),
    CreateUserResp CreateUser(1: CreateUserReq req),
}
