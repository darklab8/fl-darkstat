package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/darklab8/fl-darkstat/darkgrpc"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
)

var (
	// darkgrpc.dd84ai.com
	// 37.27.207.42:50051
	addr = flag.String("addr", "darkgrpc.dd84ai.com:80", "the address to connect to")
)

func main() {
	flag.Parse()

	c := darkgrpc.NewClient(*addr)
	defer c.Conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetBases(ctx, &pb.Empty{})
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
