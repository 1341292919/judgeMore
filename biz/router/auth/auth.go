package auth

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/mw/jwt"
	"judgeMore/biz/pack"
	"judgeMore/pkg/errno"
)

func Auth(rank int) []app.HandlerFunc {
	//为了有扩展性
	return append(make([]app.HandlerFunc, 0),
		DoubleTokenAuthFunc(rank),
	)
}

func DoubleTokenAuthFunc(rank int) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		Aerr := jwt.IsAccessTokenAvailable(ctx, c, rank)
		if Aerr != nil {
			if errors.Is(Aerr, errno.AuthAccessExpired) {
				Rerr := jwt.IsRefreshTokenAvailable(ctx, c, rank)
				if Rerr != nil {
					pack.SendFailResponse(c, errno.ConvertErr(Rerr))
					c.Abort()
					return
				}
				jwt.GenerateAccessToken(c)
			} else {
				pack.SendFailResponse(c, errno.ConvertErr(Aerr))
				c.Abort()
				return
			}
		}
		c.Next(ctx)
	}
}
