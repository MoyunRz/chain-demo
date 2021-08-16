package main

import (
	conf "chain-demo/cloudserver/hq-echo/config"
	"chain-demo/cloudserver/hq-echo/module/log"
	"chain-demo/cloudserver/hq-echo/router"
	registry2 "chain-demo/cloudserver/registry"
	"flag"
	"fmt"
	"github.com/asim/go-micro/v3/web"
	"os"
)

const (
	DefaultConfFilePath = "config/conf.toml"
)

var (
	confFilePath string
	cmdHelp      bool
)

// 初始化
func init() {
	ph, err := os.Getwd()
	if err != nil {
		fmt.Errorf("配置文件路径找不到")
	}
	// 配置配置文件所在位置
	flag.StringVar(&confFilePath, "c", ph+"/cloudserver/hq-echo/"+DefaultConfFilePath, "配置文件路径")
	// 赋予bool
	flag.BoolVar(&cmdHelp, "h", false, "帮助")
	flag.Parse()

}

func main() {
	if cmdHelp {
		// 输出默认值
		flag.PrintDefaults()
		return
	}
	// 日志打印
	log.Debugf("run with conf:%s", confFilePath)
	// 子域名部署
	router.RunSubdomains(confFilePath)
	RegisterRouter().Run()
}

func RegisterRouter() web.Service {
	conf := conf.Conf.CloudServer
	reg := registry2.RegConfig{
		Name:        conf.ServerName,
		Address:     conf.ServerAddr,
		HandlerType: "ECHO",
		EchoHandler: router.Echos,
	}
	//注册服务
	microService := registry2.InitRegistry(reg)
	return microService
}
