package config

import (
	"bytes"
	registry2 "chain-demo/cloudserver/registry"
	routers2 "chain-demo/cloudserver/userserver/routers"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/asim/go-micro/v3/web"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	Conf              cloudconfig // holds the global app config.
	defaultConfigFile = "/cloudserver/orderserver/config/config.toml"
)

type cloudconfig struct {
	// 应用配置
	CloudServer cloudserver
}

type cloudserver struct {
	ServerName string `toml:"server_name"`
	ServerAddr string `toml:"server_addr"`
}

func init() {
	InitConfig("")
}

// InitConfig initConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) error {
	op, err := os.Getwd()
	if err != nil {
		fmt.Errorf("找不到文件")
	}
	defaultConfigFile = op + defaultConfigFile
	if configFile == "" {
		configFile = defaultConfigFile
	}

	// Set defaults.
	Conf = cloudconfig{}
	if _, err := os.Stat(configFile); err != nil {
		return errors.New("config file err:" + err.Error())
	} else {
		log.Infof("load config from file:" + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.New("config load err:" + err.Error())
		}
		_, err = toml.Decode(string(configBytes), &Conf)
		if err != nil {
			return errors.New("config decode err:" + err.Error())
		}
	}

	// @TODO 配置检查
	log.Infof("config data:%v", Conf)

	return nil
}

func RegisterRouter() web.Service {
	conf := Conf.CloudServer
	reg := registry2.RegConfig{
		Name:       conf.ServerName,
		Address:    conf.ServerAddr,
		GinHandler: routers2.InitRouters(),
	}

	//注册服务
	microService := registry2.InitRegistry(reg)

	//// 首先，使用servers,err :=consulReg.GetService(serviceName)获取注册的服务
	//// 返回的servers是个slice
	hostAddress := registry2.GetServiceAddr("userserver")

	if len(hostAddress) <= 0 {
		fmt.Println("hostAddress is null")
	} else {
		url := "http://" + hostAddress + "/users"
		response, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte("")))
		fmt.Println(response)
	}

	return microService
}
