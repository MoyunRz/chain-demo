package db

import (
	"chain-demo/config"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Database struct {
	Type        string `json:"type"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	TablePrefix string `json:"table_prefix"`
}

var DatabaseSetting = &Database{}

var Db *gorm.DB

func init() {
	Setup()
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	// log.Infof(format, args...)
	fmt.Printf(format, args...)
}

func Setup() {
	// db = newConnection()
	var dbURI string
	var dialector gorm.Dialector

	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		conf.Conf.Database.UserName,
		conf.Conf.Database.Pwd,
		conf.Conf.Database.Host,
		conf.Conf.Database.Port,
		conf.Conf.Database.Name)

	dialector = mysql.New(mysql.Config{
		DSN:                       dbURI, // 数据源名称
		DefaultStringSize:         256,   // 字符串字段的默认大小
		DisableDatetimePrecision:  true,  // 禁用日期时间精度，这在 MySQL 5.6 之前不支持
		DontSupportRenameIndex:    true,  // 重命名索引时删除和创建，MySQL 5.7 之前不支持重命名索引，MariaDB
		DontSupportRenameColumn:   true,  // 重命名列时`更改` , MySQL 8 之前不支持重命名列，MariaDB
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Print(err.Error())
	}
	sqlDB, err := conn.DB()
	if err != nil {
		fmt.Errorf("connect db cloudserver failed.")
	}
	sqlDB.SetMaxIdleConns(conf.Conf.Database.MaxConn) // SetMaxIdleConns 设置空闲连接池中的最大连接数。
	sqlDB.SetMaxOpenConns(conf.Conf.Database.MaxOpen) // SetMaxOpenConns 设置到数据库的最大打开连接数。
	sqlDB.SetConnMaxLifetime(time.Second * 600)       // SetConnMaxLifetime 设置连接可以重用的最长时间。
	Db = conn
}

// GetDB 开放给外部获得db连接
func GetDB() *gorm.DB {

	sqlDB, err := Db.DB()
	if err != nil {
		fmt.Errorf("connect db cloudserver failed.")
		Setup()
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		Setup()
	}
	return Db
}
