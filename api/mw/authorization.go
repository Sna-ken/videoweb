package mw

import (
	"context"

	"github.com/Sna-ken/videoweb/pkg/base"
	"github.com/Sna-ken/videoweb/pkg/constants"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func Authorization() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		access_token := c.GetHeader(constants.AuthorizationHeader)

		claims, err := jwt.ValidateToken(string(access_token), constants.TypeAccessToken)
		if err != nil {
			base.ErrHTTPResp(c, err)
			c.Abort()
			return
		}
		c.Set(constants.UserIDPrefix, claims.Subject)
		c.Next(ctx)
	}
}

// RefreshAuthorization 用于刷新接口，允许 access_token 过期，但仍验证签名和其他声明。
func RefreshAuthorization() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		accessToken := c.GetHeader(constants.AuthorizationHeader)

		claims, err := jwt.ValidateExpiredAccessToken(string(accessToken))
		if err != nil {
			base.ErrHTTPResp(c, err)
			c.Abort()
			return
		}

		c.Set(constants.UserIDPrefix, claims.Subject)
		c.Next(ctx)
	}
}
