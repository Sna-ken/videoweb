namespace go auth

include "model.thrift"

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
    1:string id
    2: string refresh_token,
}

struct LogoutResp{
    1: model.BaseResp base,
}

struct GetMFAqrReq{
    1:string id
}

struct GetMFAqrResp{
    1: model.BaseResp base,
    2: string qr_code,
}

struct BindMFAReq{
    1:string id
    2: string mfa_code,
}

struct BindMFAResp{
    1: model.BaseResp base,
}

struct UnbindMFAReq{
    1:string id
    2: string mfa_code,
}

struct UnbindMFAResp{
    1: model.BaseResp base,
}

struct RefreshTokenReq{
    1: string id,
    2: string refresh_token,
}

struct RefreshTokenResp{
    1: model.BaseResp base,
    2: string access_token,
}

service AuthService{
    RegisterResp Register(1: RegisterReq req),
    LoginResp Login(1: LoginReq req),
    LogoutResp Logout(1: LogoutReq req),
    GetMFAqrResp GetMFAqr(1: GetMFAqrReq req),
    BindMFAResp BindMFA(1: BindMFAReq req),
    UnbindMFAResp UnbindMFA(1: UnbindMFAReq req),
    RefreshTokenResp RefreshToken(1: RefreshTokenReq req),
}
