package auth

import (
	config "chain-demo/config"
	"chain-demo/contextx"
	wd "chain-demo/middleware"
	"chain-demo/module/auth"
	"chain-demo/module/errors"
	"chain-demo/module/ginx"
	"github.com/LyricTian/gin-admin/v8/pkg/logger"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func wrapUserAuthContext(c *gin.Context, userID uint64, userName string) {
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	ctx = contextx.NewUserName(ctx, userName)
	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewUserNameContext(ctx, userName)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...wd.SkipperFunc) gin.HandlerFunc {
	if !config.Conf.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.Conf.Root.UserID, config.Conf.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if wd.SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		tokenUserID, err := a.ParseUserID(c.Request.Context(), ginx.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.Conf.IsDebugMode() {
					wrapUserAuthContext(c, config.Conf.Root.UserID, config.Conf.Root.UserName)
					c.Next()
					return
				}
				ginx.ResError(c, errors.ErrInvalidToken)
				return
			}
			ginx.ResError(c, errors.WithStack(err))
			return
		}

		idx := strings.Index(tokenUserID, "-")
		if idx == -1 {
			ginx.ResError(c, errors.ErrInvalidToken)
			return
		}

		userID, _ := strconv.ParseUint(tokenUserID[:idx], 10, 64)
		wrapUserAuthContext(c, userID, tokenUserID[idx+1:])
		c.Next()
	}
}
