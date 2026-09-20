namespace go api

include "model.thrift"

// auth部分
struct RegisterReq{
    1: string username,
    2: string password,
}

struct RegisterResp{
    1: model.BaseResp base,
}

struct LoginReq{
    1: string username,
    2: string password,
    3: string mfa_code,
}

struct LoginResp{
    1: model.BaseResp base,
    2: string access_token,
    3: string refresh_token,
}

struct LogoutReq{
}

struct LogoutResp{
    1: model.BaseResp base,
}

struct GetMFAqrReq{
}

struct GetMFAqrResp{
    1: model.BaseResp base,
    2: string qr_code,
}

struct BindMFAReq{
    1: string mfa_code,
}

struct BindMFAResp{
    1: model.BaseResp base,
}

struct UnbindMFAReq{
    1: string mfa_code,
}

struct UnbindMFAResp{
    1: model.BaseResp base,
}

struct RefreshTokenReq{
    1: string refresh_token,
}

struct RefreshTokenResp{
    1: model.BaseResp base,
    2: string access_token,
}

service AuthService{
    RegisterResp Register(1: RegisterReq req)(api.post="/api/v1/auth/register"),
    LoginResp Login(1: LoginReq req)(api.post="/api/v1/auth/login"),
    LogoutResp Logout(1: LogoutReq req)(api.post="/api/v1/auth/logout"),
    GetMFAqrResp GetMFAqr(1: GetMFAqrReq req)(api.get="/api/v1/auth/mfa/qr"),
    BindMFAResp BindMFA(1: BindMFAReq req)(api.post="/api/v1/auth/mfa/bind"),
    UnbindMFAResp UnbindMFA(1: UnbindMFAReq req)(api.post="/api/v1/auth/mfa/unbind"),
    RefreshTokenResp RefreshToken(1: RefreshTokenReq req)(api.post="/api/v1/auth/token/refresh"),
}

// user部分
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

service UserService{
    GetUserInfoResp GetUserInfo(1: GetUserInfoReq req)(api.get="/api/v1/user/info"),
    UploadAvatarResp UploadAvatar(1: UploadAvatarReq req)(api.post="/api/v1/user/avatar"),
    UpdateSignatureResp UpdateSignature(1: UpdateSignatureReq req)(api.put="/api/v1/user/signature"),
}
