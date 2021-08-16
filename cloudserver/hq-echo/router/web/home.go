package web

import (
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"

	. "chain-demo/cloudserver/hq-echo/config"
	"chain-demo/cloudserver/hq-echo/model"
	"chain-demo/cloudserver/hq-echo/model/orm"
	"chain-demo/cloudserver/hq-echo/module/log"
)

func HomeHandler(c *Context) error {
	// OpenTracing层级监控示例，API层通过中间件已支持
	span := c.OpenTracingSpan()
	if span != nil {
		// Since we have to inject our span into the HTTP headers, we create a request
		asyncReq, _ := http.NewRequest("GET", "http://"+Conf.Server.DomainApi+"/login", nil)
		// Inject the span context into the header
		err := span.Tracer().Inject(span.Context(),
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(asyncReq.Header))
		if err != nil {
			log.Debugf("Could not inject span context into header: %v", err)
		}
		go func() {
			if _, err := http.DefaultClient.Do(asyncReq); err != nil {
				span.SetTag("error", true)
				span.LogEvent(fmt.Sprintf("GET /login error: %v", err))
			}
		}()
	} else {
		log.Debugf("opentracing span nil")
	}

	User := model.User{
		Model: orm.Model{Context: c},
		Id:    1,
	}
	User.TraceGetUserById(1)

	c.Set("tmpl", "web/home")
	c.Set("data", map[string]interface{}{
		"title": "Home",
	})

	return nil
}
