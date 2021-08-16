package main

import (
	"chain-demo/cloudserver/gateway/proto/hello"
	"context"
	"fmt"
	"github.com/micro/micro/v3/service"
	"time"
)

func main() {
	srv := service.New()
	// create the proto client for helloworld
	client := hello.NewHelloService("hello", srv.Client())
	// call an endpoint on the service
	rsp, err := client.Call(context.Background(), &hello.Request{
		Name: "John",
	})

	if err != nil {
		fmt.Println("Error calling hello: ", err)
		return
	}
	// 打印响应内容
	fmt.Println("Response: ", rsp.Msg)
	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}
