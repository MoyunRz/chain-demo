package config

import (
	registry2 "chain-demo/cloudserver/registry"
	"chain-demo/cloudserver/userserver/routers"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/asim/go-micro/v3/web"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
)

var (
	Conf              cloudconfig // holds the global app config.
	defaultConfigFile = "/cloudserver/userserver/config/config.toml"
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
		Name:        conf.ServerName,
		Address:     conf.ServerAddr,
		HandlerType: "GIN",
		GinHandler:  routers.InitRouters(),
	}
	microService := registry2.InitRegistry(reg)

	return microService
}
