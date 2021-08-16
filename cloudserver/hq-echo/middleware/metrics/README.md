# Metrics

## Grafana
浏览器:localhost:3000
```bash
docker run --name grafana -d -p 3000:3000 grafana/grafana
```

## Prometheus
浏览器:localhost:9090
```bash
docker run -d --name prometheus -p 9090:9090 -v ~/tmp/prometheus.yml:/etc/prometheus/prometheus.yml \
              prom/prometheus
              
# Dashboard JSON
# middleware/metrics/prometheus/grafana.json
```

配置文件~/tmp/prometheus.yml
> Docker内网络需要能访问到域名:echo.www.localhost.com
> 可以在/etc/hosts将域名映射到主机IP，可以参考README.md [自动修改/etch/osts映射自定义域名到主机ip](/README.md#%E8%87%AA%E5%8A%A8%E4%BF%AE%E6%94%B9etchosts%E6%98%A0%E5%B0%84%E8%87%AA%E5%AE%9A%E4%B9%89%E5%9F%9F%E5%90%8D%E5%88%B0%E4%B8%BB%E6%9C%BAip)
```bash
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # Evaluate rules every 15 seconds.

  # Attach these extra labels to all timeseries collected by this Prometheus instance.
  external_labels:
    monitor: 'codelab-monitor'

rule_files:
  - 'prometheus.rules.yml'

scrape_configs:
  - job_name: 'prometheus'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']

  - job_name:       'echo-web'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['echo.www.localhost.com']
        labels:
          group: 'production'

```

## Push模式
### Graphite
浏览器:localhost:8090
```sh
# 登录账户名:root，密码:root
docker run -d --name graphite --restart=always -p 8090:80 -p 2003-2004:2003-2004 -p 2023-2024:2023-2024 -p 8125:8125/udp -p 8126:8126 hopsoft/graphite-statsd

# Dashboard JSON
# middleware/metrics/grafana_graphite.json
```

### InfluxDB
```bash
# 未测试，省略……😂
```
