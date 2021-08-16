package registry

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
	"github.com/asim/go-micro/v3/web"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var consulReg registry.Registry

type RegConfig struct {
	Name        string
	Address     string
	HandlerType string
	GinHandler  *gin.Engine
	EchoHandler *echo.Echo
}

func init() {
	//新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
	consulReg = consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)
}

func MatchType(reg RegConfig) http.Handler {

	switch reg.HandlerType {
	case "GIN":
		return reg.GinHandler
	case "ECHO":
		return reg.EchoHandler
	default:
		return reg.GinHandler
	}
}

func InitRegistry(reg RegConfig) web.Service {

	//注册服务
	microService := web.NewService(
		web.Name(reg.Name),
		//web.RegisterTTL(time.Second*30),//设置注册服务的过期时间
		//web.RegisterInterval(time.Second*20),//设置间隔多久再次注册服务
		web.Address(reg.Address),
		web.Handler(MatchType(reg)),
		web.Registry(consulReg),
	)

	return microService
}

func GetServiceAddr(serviceName string) (address string) {
	var retryCount int
	for {
		servers, err := consulReg.GetService(serviceName)
		if err != nil {
			fmt.Println(err.Error())
		}
		var services []*registry.Service
		for _, value := range servers {
			fmt.Println(value.Name, ":", value.Version)
			services = append(services, value)
		}
		// 使用next := selector.RoundRobin(services)获取其中一个服务的信息
		next := selector.RoundRobin(services)
		if node, err := next(); err == nil {
			address = node.Address
		}
		if len(address) > 0 {
			return
		}
		//重试次数++
		retryCount++
		time.Sleep(time.Second * 1)
		//重试5次为获取返回空
		if retryCount >= 5 {
			return
		}
	}
}
