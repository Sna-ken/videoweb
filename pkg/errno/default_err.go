package errno

var (
	Success               = NewErr(SuccessCode, "成功")
	UnknownError          = NewErr(UnknownErrorCode, "未知错误")
	ParamError            = NewErr(ParamErrorCode, "参数错误")
	ParamEmptyError       = NewErr(ParamEmptyErrorCode, "参数为空")
	ParamTypeError        = NewErr(ParamTypeErrorCode, "参数类型错误")
	ParamInvalidError     = NewErr(ParamInvalidErrorCode, "参数无效")
	AuthError             = NewErr(AuthErrorCode, "认证错误")
	AuthExpiredError      = NewErr(AuthExpiredErrorCode, "认证过期")
	AuthInvalidError      = NewErr(AuthInvalidErrorCode, "认证无效")
	InternalServiceError  = NewErr(InternalServiceErrorCode, "内部服务错误")
	InternalDatabaseError = NewErr(InternalDatabaseErrorCode, "内部数据库错误")
)
