package session

import (
	"github.com/labstack/echo/v4"

	. "chain-demo/cloudserver/hq-echo/config"
	es "chain-demo/cloudserver/hq-echo/middleware/session"
)

func Session() echo.MiddlewareFunc {
	// 匹配缓存的参数，new出相应的组件
	switch Conf.SessionStore {
	case REDIS:
		// redis：key-value
		store, err := es.NewRedisStore(10, "tcp", Conf.Redis.Server, Conf.Redis.Pwd, []byte("secret-key"))
		if err != nil {
			panic(err)
		}
		return es.New("sid", store)
	case FILE:
		// 文件类型存储
		store := es.NewFilesystemStore("", []byte("secret-key"))
		return es.New("sid", store)
	default:
		// 默认cookie存储
		store := es.NewCookieStore([]byte("secret-key"))
		return es.New("sid", store)
	}
}
