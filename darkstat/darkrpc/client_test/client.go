package main

import (
	"fmt"
	"log"

	"github.com/darklab8/fl-darkstat/darkstat/darkrpc"
)

func GetBases(args darkrpc.Args, reply *darkrpc.Reply) error {
	return nil
}

func main() {
	smth := darkrpc.ServerRpc{}
	smth2 := smth.GetBases
	_ = smth2
	smth2 = GetBases

	args := darkrpc.Args{}
	var reply darkrpc.Reply

	client := darkrpc.NewClient()

	err := client.GetBases(args, &reply)
	if err != nil {
		log.Fatal("getBases error:", err)
	}
	fmt.Println("Bases[0]=", reply.Bases[0])

}
