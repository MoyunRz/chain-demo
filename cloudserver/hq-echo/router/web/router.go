package web

import (
	. "chain-demo/cloudserver/hq-echo/config"
	"chain-demo/cloudserver/hq-echo/middleware/captcha"
	"chain-demo/cloudserver/hq-echo/middleware/opentracing"
	"chain-demo/cloudserver/hq-echo/model"
	"chain-demo/cloudserver/hq-echo/module/auth"
	"chain-demo/cloudserver/hq-echo/module/cache"
	"chain-demo/cloudserver/hq-echo/module/render"
	"chain-demo/cloudserver/hq-echo/module/session"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

// Routers ---------
// Website Routers
// ---------
func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Context自定义
	e.Use(NewContext())

	// Customization
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("web")
	e.Logger.SetLevel(GetLogLvl())

	//// 静态资源
	//switch Conf.Static.Type {
	//case BINDATA:
	//	e.Use(staticbin.Static(assets.Asset, staticbin.Options{
	//		Dir:         "/",
	//		SkipLogging: true,
	//	}))
	//default:
	//	e.Static("/assets", "./assets")
	//}

	// Session
	e.Use(session.Session())

	// CSRF
	e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
		ContextKey:  "_csrf",
		TokenLookup: "form:_csrf",
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// 验证码，优先于静态资源
	e.Use(captcha.Captcha(captcha.Config{
		CaptchaPath: "/captcha/",
		SkipLogging: true,
	}))

	// Gzip，在验证码、静态资源之后
	// 验证码、静态资源使用http.ServeContent()，与Gzip有冲突，Nginx报错，验证码无法访问
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// OpenTracing
	if !Conf.Opentracing.Disable {
		e.Use(opentracing.OpenTracing("web"))
	}

	// 模板
	e.Renderer = render.LoadTemplates()
	e.Use(render.Render())

	// Cache
	e.Use(cache.Cache())

	// Auth
	e.Use(auth.New(model.GenerateAnonymousUser))
	//e.Use(rabbitmq.NewMQQueue(rabbitmq.DefaultMQConfig))
	// Routers
	e.GET("/", handler(HomeHandler))
	e.GET("/login", handler(LoginHandler))
	e.GET("/register", handler(RegisterHandler))
	e.GET("/logout", handler(LogoutHandler))
	e.POST("/login", handler(LoginPostHandler))
	e.POST("/register", handler(RegisterPostHandler))

	e.GET("/dashboard", DashboardHandler)

	//e.GET("/jwt/tester", handler(JWTTesterHandler))
	e.GET("/ws", handler(WsHandler))

	user := e.Group("/user")
	user.Use(auth.LoginRequired())
	{
		user.GET("/:id", handler(UserHandler))
	}

	about := e.Group("/about")
	about.Use(auth.LoginRequired())
	{
		about.GET("", handler(AboutHandler))
	}

	return e
}

type (
	HandlerFunc func(*Context) error
)

/**
 * 自定义Context的Handler
 */
func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
