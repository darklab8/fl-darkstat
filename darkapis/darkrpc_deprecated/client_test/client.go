package main

import (
	"log"

	"github.com/darklab8/fl-darkstat/darkapis/darkrpc_deprecated"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
)

func GetBases(args darkrpc_deprecated.Args, reply *darkrpc_deprecated.Reply) error {
	return nil
}

func main() {
	smth := darkrpc_deprecated.ServerRpc{}
	smth2 := smth.GetBases
	_ = smth2

	args := darkrpc_deprecated.Args{}
	var reply darkrpc_deprecated.Reply

	client := darkrpc_deprecated.NewClient(darkrpc_deprecated.WithPortCli(8111))
	// client := darkrpc.NewClient(darkrpc.WithSockCli(darkrpc.DarkstatRpcSock))

	var health_reply bool
	logus.LogCli.Infoln("attempted to get health")
	err := client.GetHealth(args, &health_reply)
	if err != nil {
		log.Fatal("getHealth error:", err)
	} else {
		logus.LogCli.Infoln("server health reply =", health_reply)
	}

	err = client.GetBases(args, &reply)
	if err != nil {
		log.Fatal("getBases error:", err)
	}
	logus.LogCli.Infoln("Bases[0]=", reply.Bases[0])
}
