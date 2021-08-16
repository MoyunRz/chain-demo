// Package main
package main

import (
	"context"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/registry"
	"time"

	hello "chain-demo/cloudserver/greeter/srv/proto/hello"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/util/log"
	"google.golang.org/grpc"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Log("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	go func() {
		for {
			grpc.DialContext(context.TODO(), "127.0.0.1:9091")
			time.Sleep(time.Second)
		}
	}()
	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)
	srv := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.Address(":18005"),
		micro.Registry(consulReg),
	)
	// optionally setup command line usage
	srv.Init()
	// Register Handlers
	hello.RegisterSayHandler(srv.Server(), new(Say))
	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
