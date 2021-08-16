package main

import (
	"chain-demo/cloudserver/userserver/config"
)

func main() {
	config.RegisterRouter().Run()
}
