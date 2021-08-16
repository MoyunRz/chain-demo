package orm

import (
	"chain-demo/cloudserver/hq-echo/module/log"
)

type Logger struct {
}

// Print format & print log
func (logger Logger) Print(values ...interface{}) {
	// @TODO
	// 日志格式化解析
	log.Debugf("orm log:%v \n", values)
}
