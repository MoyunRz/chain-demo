package router

import (
	. "chain-demo/cloudserver/hq-echo/config"
	"chain-demo/cloudserver/hq-echo/middleware/metrics/prometheus"
	"chain-demo/cloudserver/hq-echo/middleware/opentracing"
	"chain-demo/cloudserver/hq-echo/middleware/pprof"
	"chain-demo/cloudserver/hq-echo/module/log"
	"chain-demo/cloudserver/hq-echo/router/web"
	"context"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type (
	// Host 结构体
	Host struct {
		Echo *echo.Echo
	}
)

var Echos *echo.Echo

// InitRoutes 初始化路由
func InitRoutes() map[string]*Host {
	// Hosts
	hosts := make(map[string]*Host) /*创建集合 */
	// Key-value
	hosts[Conf.Server.DomainWeb] = &Host{web.Routers()}
	//hosts[Conf.Server.DomainSocket] = &Host{socket.Routers()}
	// 路由集合
	return hosts
}

// RunSubdomains 子域名部署
func RunSubdomains(confFilePath string) {
	// 配置初始化
	if err := InitConfig(confFilePath); err != nil {
		log.Panic(err)
	}

	// 全局日志级别
	log.SetLevel(GetLogLvl())

	// Server
	Echos = echo.New()

	// pprof
	Echos.Pre(pprof.Serve())

	Echos.Pre(mw.RemoveTrailingSlash())

	// 请求追踪
	// Elastic APM
	// Requires APM Server 6.5.0 or newer
	apm.DefaultTracer.Service.Name = Conf.Opentracing.ServiceName
	apm.DefaultTracer.Service.Version = Conf.App.Version
	Echos.Use(apmechov4.Middleware(
		apmechov4.WithRequestIgnorer(func(request *http.Request) bool {
			return false
		}),
	))

	// OpenTracing
	otCtf := opentracing.Configuration{
		Disabled: Conf.Opentracing.Disable,
		Type:     opentracing.TracerType(Conf.Opentracing.Type),
	}
	if closer := otCtf.InitGlobalTracer(
		opentracing.ServiceName(Conf.Opentracing.ServiceName),
		opentracing.Address(Conf.Opentracing.Address),
	); closer != nil {
		defer closer.Close()
	}

	// 日志级别
	Echos.Logger.SetLevel(GetLogLvl())

	// Metrics
	if !Conf.Metrics.Disable {
		Echos.Use(prometheus.MetricsFunc(
			prometheus.Namespace("echo_web"),
		))

		/* Push模式
		m := metrics.NewMetrics(metrics.Prefix(""))
		echos.Use(metrics.MetricsFunc(m))
		m.MemStats()

		hostname, err := os.Hostname()
		if err != nil {
			log.Warnf("os.Hostname() error:%v", err)
			hostname = "-"
		}
		// Graphite
		addr, _ := net.ResolveTCPAddr("tcp", Conf.Metrics.Address)
		m.Graphite(Conf.Metrics.FreqSec*time.Second, "echo-web.node."+hostname, addr)

		// InfluxDB
		m.InfluxDBWithTags(
			Conf.Metrics.FreqSec*time.Second,
			"http://127.0.0.1:8086",
			"metrics",
			"test",
			"test",
			map[string]string{"node": hostname})
		*/
	}

	// Secure, XSS/CSS HSTS
	Echos.Use(mw.SecureWithConfig(mw.DefaultSecureConfig))
	Echos.Use(mw.MethodOverride())

	// CORS
	Echos.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins: []string{"http://" + Conf.Server.DomainWeb, "http://" + Conf.Server.DomainApi},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization},
	}))
	Echos.Use()
	hosts := InitRoutes()
	Echos.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		// "postgres://user:pass@host.com:5432/path?k=v#f"
		// 我们将解析这个 URL 示例，它包含了一个 scheme，认证信息，主机名，端口，路径，查询参数和片段。
		u, _err := url.Parse(c.Scheme() + "://" + req.Host)
		if _err != nil {
			Echos.Logger.Errorf("Request URL parse error:%v", _err)
		}

		host := hosts[u.Hostname()]
		if host == nil {
			Echos.Logger.Info("Host not found")
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return
	})

	//if !Conf.Server.Graceful {
	//	Echos.Logger.Fatal(Echos.Start(Conf.Server.Addr))
	//} else {
	if false {
		// Graceful Shutdown
		// Start server
		go func() {
			if err := Echos.Start(Conf.Server.Addr); err != nil {
				Echos.Logger.Errorf("Shutting down the server with error:%v", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := Echos.Shutdown(ctx); err != nil {
			Echos.Logger.Fatal(err)
		}
	}
}
