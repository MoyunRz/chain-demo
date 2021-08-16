package log

/**
 *
 * 日志文件
 *
 */

import (
	l "github.com/labstack/gommon/log"
)

// Logger 结构体
type Logger struct {
	*l.Logger // 输出Logger的结构体的值
}

var (
	// 自定义全局log
	// 使用gommon/log的global由于限制3层调用栈获取不到log的准确file路径
	global = l.New("log")
	// 用于日志标头的设置，默认值为
	defaultHeader = `{"time":"${time_rfc3339}","level":"${level}","prefix":"${prefix}",` +
		`"file":"${long_file}","line":"${line}"}`
)

// init 初始化
func init() {
	// 全局log属性设置
	l.SetHeader(defaultHeader)
	l.SetLevel(l.DEBUG)

	// 自定义全局logs属性设置
	global.SetHeader(defaultHeader)
	global.SetLevel(l.DEBUG)
}

func SetLevel(v l.Lvl) {
	l.SetLevel(v)
	global.SetLevel(v)
}

func Debug(i ...interface{}) {
	global.Debug(i)
}

func Debugf(format string, values ...interface{}) {
	global.Debugf(format, values...)
}

func Info(i ...interface{}) {
	global.Info(i)
}

func Infof(format string, values ...interface{}) {
	global.Infof(format, values...)
}

func Warn(i ...interface{}) {
	global.Warn(i)
}

func Warnf(format string, values ...interface{}) {
	global.Warnf(format, values...)
}

func Error(i ...interface{}) {
	global.Error(i)
}

func Errorf(format string, values ...interface{}) {
	global.Errorf(format, values...)
}

func Fatal(i ...interface{}) {
	global.Fatal(i)
}

func Fatalf(format string, values ...interface{}) {
	global.Fatalf(format, values...)
}

func Panic(i ...interface{}) {
	global.Panic(i)
}

func Panicf(format string, args ...interface{}) {
	global.Panicf(format, args)
}
