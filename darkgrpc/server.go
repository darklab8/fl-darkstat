package darkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedDarkGRpcServer
	app_data *appdata.AppData
	port     int
}

func NewServer(app_data *appdata.AppData, port int) *Server {
	return &Server{app_data: app_data, port: port}
}

const DefaultServerPort = 50051

func (r *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDarkGRpcServer(s, r)
	log.Printf("server listening at %v", lis.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("grpc server is launched")
}

func (s *Server) GetHealth(_ context.Context, in *pb.Empty) (*pb.HealthReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	return &pb.HealthReply{
		IsHealthy: true,
	}, nil
}
