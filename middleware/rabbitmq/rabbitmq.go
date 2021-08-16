package rabbitmq

import (
	"bytes"
	"chain-demo/config"
	"github.com/gin-gonic/gin"
)

type mqconf struct {
	mtype        string
	mqueue       string
	key          string
	exchangeName string
}

const (
	WORK = "WORK"

	SIMPLE = "SIMPLE"

	PUBLISH = "PUBLISH"

	ROUTING = "ROUTING"

	TOPIC = "TOPIC"
)

var murl = ""

func init() {
	var buffer bytes.Buffer
	buffer.WriteString("amqp://")
	buffer.WriteString(conf.Conf.RabbitMQ.MQUser)
	buffer.WriteString(":")
	buffer.WriteString(conf.Conf.RabbitMQ.MQPassword)
	buffer.WriteString("@")
	buffer.WriteString(conf.Conf.RabbitMQ.MQUrl)
	murl = buffer.String()
}

var (
	// DefaultMQConfig is the default Secure middleware config.
	DefaultMQConfig = mqconf{
		mtype:  WORK,
		mqueue: "default_simple",
	}
)

// MQMiddleware 初始化MQ，并发送msg
func MQMiddleware(m mqconf, msg string) gin.HandlerFunc {

	return func(c *gin.Context) {
		switch m.mtype {
		case WORK:
			rmq := NewRabbitMQSimple(m.mqueue, "", "", murl)
			rmq.PublishSimple(msg)
		case SIMPLE:
			rmq := NewRabbitMQSimple(m.mqueue, "", "", murl)
			rmq.PublishSimple(msg)
		case PUBLISH:
			rmq := NewRabbitMQPubSub(m.mqueue, m.exchangeName, m.key, murl)
			rmq.PublishPub(msg)
		case ROUTING:
			rmq := NewRabbitMQRouting(m.mqueue, m.exchangeName, m.key, murl)
			rmq.PublishRouting(msg)
		case TOPIC:
			rmq := NewRabbitMQTopic(m.mqueue, m.exchangeName, m.key, murl)
			rmq.PublishTopic(msg)
		}
		c.Next()
	}
}
