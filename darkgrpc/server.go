package darkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedDarkGRpcServer
	app_data     *appdata.AppData
	port         int
	sock_address string
}

func NewServer(app_data *appdata.AppData, port int, opts ...ServerOpt) *Server {
	s := &Server{app_data: app_data, port: port}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithSockAddr(sock string) ServerOpt {
	return func(s *Server) {
		s.sock_address = sock
	}
}

type ServerOpt func(s *Server)

const DefaultServerPort = 50051

const DarkstatGRpcSock = "/tmp/darkstat/grpc.sock"

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

	if cfg.IsLinux && r.sock_address != "" {
		os.Remove(r.sock_address)
		os.Mkdir("/tmp/darkstat", 0777)
		sock_listener, err := net.Listen("unix", fmt.Sprintf("%s", r.sock_address)) // if serving over Unix
		if err != nil {
			log.Fatal("listen error:", err)
		}
		fmt.Println("turning on server")
		go func() {
			err = s.Serve(sock_listener)
			logus.Log.CheckError(err, "failed to serve grpc sock")
		}()
	}

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
