package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
)

var (
	// darkgrpc.dd84ai.com
	// 37.27.207.42:50051
	// fmt.Sprintf("unix:%s", darkgrpc.DarkstatGRpcSock)
	// darkgrpc-staging.dd84ai.com:443
	addr = flag.String("addr", "darkgrpc-staging.dd84ai.com:443", "the address to connect to")
)

func main() {
	flag.Parse()

	c := darkgrpc.NewClient(*addr)
	defer func() {
		err := c.Conn.Close()
		logus.Log.CheckError(err, "failed to close connection")
	}()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetBasesNpc(ctx, &pb.GetBasesInput{})
	fmt.Println(r)
	fmt.Println("err=", err)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Items[0])
	if len(r.Items) > 0 {
		fmt.Println("SUCCCCCCESSS")
	}
}
