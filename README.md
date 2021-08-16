# chain-demo

该框架是 micro为基础作为服务注册，集成gin和echo框架

用consul作为服务发现

在启动微服务后可以使用micro-api作为网关

在该项目中也集成了grpc的server和client案例

目录简介：

- base              公共基础model
- cloudserver       微服务模块
- config            配置文件
- contextx          上下文处理
- middleware        中间件
- module            公共基础模块
- proto             grpc配置模块


**cloudserver**

- gateway           网关层grpc
- greeter           micro的官方Greeter案例
- hq-echo           echo模块
- orderserver       基于http的服务发现
- registry          服务注册
- userserver        用户服务（demo）