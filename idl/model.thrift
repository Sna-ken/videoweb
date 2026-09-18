namespace go model

struct BaseResp{
    1: i64 code,
    2: string msg,
}

struct UserInfo{
    1: string user_id,
    2: string username,
    3: string avatar,
    4: string signature,
    5: i32 follow_count,
    6: i32 follower_count,
}