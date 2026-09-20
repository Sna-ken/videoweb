package errno

const (
	SuccessCode = 0

	// 公共错误 100000-100100
	UnknownErrorCode          = 100000 //未知错误
	ParamErrorCode            = 100001 //参数错误
	ParamEmptyErrorCode       = 100002 //参数为空
	ParamTypeErrorCode        = 100003 //参数类型错误
	ParamInvalidErrorCode     = 100004 //参数无效
	AuthErrorCode             = 100005 //认证错误
	AuthExpiredErrorCode      = 100006 //认证过期
	AuthInvalidErrorCode      = 100007 //认证无效
	InternalServiceErrorCode  = 100008 //内部服务错误
	InternalDatabaseErrorCode = 100009 //内部数据库错误

	// auth错误 100101-100200
	MFAcodeErrorCode           = 100101 //MFA验证码错误
	MFAcodeExpiredErrorCode    = 100102 //MFA验证码过期
	MFAcodeEmptyErrorCode      = 100103 //MFA验证码为空
	MFANotEnabledErrorCode     = 100104 //MFA未启用
	PasswordIncorrectErrorCode = 100105 //密码错误

	// user错误 100201-100300
	UserNotFoundErrorCode    = 100201 //用户不存在
	UserHasExistedErrorCode  = 100202 //用户已存在
	PasswordErrorCode        = 100203 //密码错误
	PasswordHashErrorCode    = 100204 //密码哈希错误
	PasswordInvalidErrorCode = 100205 //密码无效

)
