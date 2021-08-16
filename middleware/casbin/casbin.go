package casbin

import (
	config "chain-demo/config"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// MysqlTool *gorm.DB
func MysqlTool() *gorm.DB {
	//var err error

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Conf.Database.UserName,
		config.Conf.Database.Pwd,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.Name)

	mysql, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connect DB error")
		panic(err)
	}
	return mysql
}

func EnforcerTool() *casbin.Enforcer {
	adapter := gormadapter.NewAdapterByDB(MysqlTool())
	//创建一个casbin决策器需要一个模板文件和策略文件为参数
	Enforcer := casbin.NewEnforcer("config/keymatch.conf", adapter)
	//从数据库加载策略
	Enforcer.LoadPolicy()
	return Enforcer
}

func interceptor(e *casbin.Enforcer) gin.HandlerFunc {

	return func(c *gin.Context) {
		//获取资源
		obj := c.Request.URL.RequestURI()
		//获取方法
		act := c.Request.Method
		//获取实体
		sub := "admin"

		//判断策略中是否存在
		enforce := e.AddPolicy(sub, obj, act)
		if enforce {
			fmt.Println("通过权限")
			c.Next()
		} else {
			fmt.Println("没有通过权限")
			c.Abort()
		}

	}
}
