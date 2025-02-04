package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/darklab8/fl-darkstat/darkgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// darkgrpc.dd84ai.com
	// 37.27.207.42:50051
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDarkGRpcClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetBases(ctx, &pb.Empty{})
	fmt.Println(r)
	fmt.Println("err=", err)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Bases[0])
}
