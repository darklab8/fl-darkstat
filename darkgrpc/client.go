package darkgrpc

import (
	"crypto/tls"
	"log"
	"strings"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn *grpc.ClientConn
	pb.DarkstatClient
}

func NewClient(address string) *Client {
	var creds credentials.TransportCredentials
	if strings.Contains(address, "443") {
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: false}) // for darkstat.dd84ai.com:443
	} else {
		creds = insecure.NewCredentials() // for darkstat.dd84ai.com:80
	}

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32*10e6)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewDarkstatClient(conn)

	return &Client{
		DarkstatClient: c,
		Conn:           conn,
	}
}
