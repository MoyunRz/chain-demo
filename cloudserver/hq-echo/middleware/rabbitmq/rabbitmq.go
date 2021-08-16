package rabbitmq

import (
	"bytes"
	. "chain-demo/cloudserver/hq-echo/config"
	"github.com/labstack/echo/v4"
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
	buffer.WriteString(Conf.RabbitMQ.MQUser)
	buffer.WriteString(":")
	buffer.WriteString(Conf.RabbitMQ.MQPassword)
	buffer.WriteString("@")
	buffer.WriteString(Conf.RabbitMQ.MQUrl)
	murl = buffer.String()
}

var (
	// DefaultMQConfig is the default Secure middleware config.
	DefaultMQConfig = mqconf{
		mtype:  WORK,
		mqueue: "default_simple",
	}
)

func NewMQQueue(m mqconf) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch m.mtype {
			case WORK:
				rmq := NewRabbitMQSimple(m.mqueue, "", "", murl)
				rmq.PublishSimple(c.QueryString())
			case SIMPLE:
				rmq := NewRabbitMQSimple(m.mqueue, "", "", murl)
				rmq.PublishSimple(c.QueryString())
			case PUBLISH:
				rmq := NewRabbitMQPubSub(m.mqueue, m.exchangeName, m.key, murl)
				rmq.PublishPub(c.QueryString())
			case ROUTING:
				rmq := NewRabbitMQRouting(m.mqueue, m.exchangeName, m.key, murl)
				rmq.PublishRouting(c.QueryString())
			case TOPIC:
				rmq := NewRabbitMQTopic(m.mqueue, m.exchangeName, m.key, murl)
				rmq.PublishTopic(c.QueryString())
			default:
				return nil
			}
			return next(c)
		}
	}
}
