package metrics

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rcrowley/go-metrics"
	influxdb "github.com/vrischmann/go-metrics-influxdb"
)

// MetricsFunc 监控测量方法
func MetricsFunc(m *Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			latency := stop.Sub(start)
			status := res.Status
			host := strings.Replace(req.Host, ".", "_", -1)

			//Total request count. 请求总数
			meter := metrics.GetOrRegisterMeter(m.withPrefix("request_total."+".host."+host+".status."+strconv.Itoa(status)), m.opts.Registry)
			meter.Mark(1)

			//Request size in bytes. 请求字节长度
			meter = metrics.GetOrRegisterMeter(m.withPrefix("request_size."+".host."+host+".status."+strconv.Itoa(status)), m.opts.Registry)
			meter.Mark(req.ContentLength)

			//Request duration in nanoseconds. 以纳秒为单位的请求持续时间
			h := metrics.GetOrRegisterHistogram(m.withPrefix("request_duration."+".host."+host+".status."+strconv.Itoa(status)), m.opts.Registry,
				metrics.NewExpDecaySample(1028, 0.015))
			h.Update(latency.Nanoseconds())

			//Response size in bytes.响应大小（以字节为单位）。
			meter = metrics.GetOrRegisterMeter(m.withPrefix("response_size."+".host."+host+".status."+strconv.Itoa(status)), m.opts.Registry)
			meter.Mark(res.Size)

			return nil
		}
	}
}

type Metrics struct {
	opts Options
}

//NewMetrics 创建一个新的监控
func NewMetrics(options ...Option) *Metrics {
	opts := applyOptions(options...)
	return &Metrics{opts: opts}
}

func (m *Metrics) withPrefix(s string) string {
	return m.opts.Prefix + "." + s
}

func (m *Metrics) MemStats() {
	metrics.RegisterRuntimeMemStats(m.opts.Registry)
	go metrics.CaptureRuntimeMemStats(m.opts.Registry, 5*time.Second)
}

// Log reports metrics into logs.
//
// m.Log(30*time.Second, e.Logger)
//
func (m *Metrics) Log(freq time.Duration, l metrics.Logger) {
	go metrics.Log(m.opts.Registry, freq, l)
}

// Graphite reports metrics into graphite.
//
// 	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
//  m.Graphite(10e9, "metrics", addr)
//
func (m *Metrics) Graphite(freq time.Duration, prefix string, addr *net.TCPAddr) {
	go metrics.Graphite(m.opts.Registry, freq, prefix, addr)
}

// InfluxDB reports metrics into influxdb.
//
// 	m.InfluxDB(10e9, "http://127.0.0.1:8086","metrics", "test","test"})
//
func (m *Metrics) InfluxDB(freq time.Duration, url, database, username, password string) {
	go influxdb.InfluxDB(m.opts.Registry, freq, url, database, username, password)
}

// InfluxDBWithTags reports metrics into influxdb with tags.
// you can set node info into tags.
//
// 	m.InfluxDBWithTags(10e9, "http://127.0.0.1:8086","metrics", "test","test", map[string]string{"host":"127.0.0.1"})
//
func (m *Metrics) InfluxDBWithTags(freq time.Duration, url, database, username, password string, tags map[string]string) {
	go influxdb.InfluxDBWithTags(m.opts.Registry, freq, url, database, username, password, tags)
}
