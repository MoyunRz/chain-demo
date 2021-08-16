文件目录结构：
|---hq-echo
    |---assets      前端静态资源
    |---config      配置文件
    |---middleware  中间件
    |---model       项目模型，实体
    |---module      项目模块
    |---proto       编译配置文件
    |---router      路由管理
    |---template    模板
    |---util        工具
    |---go.mod      包版本管理
    |---main.go     项目启动

框架使用：

- go: 1.16
- echo:v4
- redis
- rabbitmq
- protoc/grpc
- mysql
- metrics
- opentracing

**grpc**

配置文件位置：proto/*

文件格:*.proto

执行指令：
```shell
  // '-I（-IPATH）指定要在其中搜索导入（import）的目录。'
  // protoc {输出目录} {proto文件位置} --go_out=plugins=grpc:{输出目录}
  protoc -I helloworld/ helloworld/pb/helloworld.proto --go_out=plugins=grpc:helloworld
  // helloworld/pb目录下会生成 helloworld.pb.go 文件
```
proto生成文件目录：module/grpc/{文件夹名}/

使用grpc案例：

```go
package main 
import (···)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	// 创建gRPC服务器
	s := grpc.NewServer()
	// 在gRPC服务端注册服务
	pb.RegisterGreeterServer(s, &server{}) 
	//在给定的gRPC服务器上注册服务器反射服务
	reflection.Register(s) 
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
```


