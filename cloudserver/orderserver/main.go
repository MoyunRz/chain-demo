package main

import (
	"chain-demo/cloudserver/orderserver/config"
)

func main() {
	config.RegisterRouter().Run()
}
