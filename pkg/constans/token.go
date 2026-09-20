package constans

import "time"

var (
	TypeAccessToken    = "access_token"
	TypeRefreshToken   = "refresh_token"
	AccessTokenExpire  = 30 * time.Minute    //30分钟
	RefreshTokenExpire = 30 * 24 * time.Hour //30天

	Issuer = "Sna-ken"

	RefreshTokenKeyPrefix = "refreshtoken:"
)
