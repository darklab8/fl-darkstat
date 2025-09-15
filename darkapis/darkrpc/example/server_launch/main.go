package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/darklab8/fl-darkstat/darkapis/darkrpc/example/server"
)

func main() {
	arith := new(server.Arith)
	_ = rpc.Register(arith)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	_ = http.Serve(l, nil)
}
