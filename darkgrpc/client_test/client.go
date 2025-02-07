package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/darklab8/fl-darkstat/darkgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// darkgrpc.dd84ai.com
	// 37.27.207.42:50051
	addr = flag.String("addr", "darkgrpc.dd84ai.com:80", "the address to connect to")
)

func main() {
	flag.Parse()
	var creds credentials.TransportCredentials
	if strings.Contains(*addr, "443") {
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: false}) // for darkstat.dd84ai.com:443
	} else {
		creds = insecure.NewCredentials() // for darkstat.dd84ai.com:80
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
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
	log.Printf("Greeting: %s", r.Items[0])
	if len(r.Items) > 0 {
		fmt.Println("SUCCCCCCESSS")
	}
}
