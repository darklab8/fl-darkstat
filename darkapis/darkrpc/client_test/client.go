package main

import (
	"fmt"
	"log"

	"github.com/darklab8/fl-darkstat/darkapis/darkrpc"
)

func GetBases(args darkrpc.Args, reply *darkrpc.Reply) error {
	return nil
}

func main() {
	smth := darkrpc.ServerRpc{}
	smth2 := smth.GetBases
	_ = smth2

	args := darkrpc.Args{}
	var reply darkrpc.Reply

	client := darkrpc.NewClient(darkrpc.WithPortCli(8111))
	// client := darkrpc.NewClient(darkrpc.WithSockCli(darkrpc.DarkstatRpcSock))

	var health_reply bool
	fmt.Println("attempted to get health")
	err := client.GetHealth(args, &health_reply)
	if err != nil {
		log.Fatal("getHealth error:", err)
	} else {
		fmt.Println("server health reply =", health_reply)
	}

	err = client.GetBases(args, &reply)
	if err != nil {
		log.Fatal("getBases error:", err)
	}
	fmt.Println("Bases[0]=", reply.Bases[0])
}
