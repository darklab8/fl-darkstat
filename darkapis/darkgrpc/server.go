package darkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc/staticproto"
	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	_ "github.com/darklab8/fl-darkstat/docs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedDarkstatServer
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
	s := grpc.NewServer(grpc.MaxRecvMsgSize(32 * 10e6))
	pb.RegisterDarkstatServer(s, r)
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

	{
		// GRPC GATEWAY https://github.com/grpc-ecosystem/grpc-gateway
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32 * 10e6)),
		}

		var err error
		if cfg.IsLinux {
			err = pb.RegisterDarkstatHandlerFromEndpoint(ctx, mux, fmt.Sprintf("unix:%s", DarkstatGRpcSock), opts)
		} else {
			err = pb.RegisterDarkstatHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
		}
		if err != nil {

			panic(err)
		}

		// Sprinked with API documentation :)
		mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte(staticproto.Index))
		})
		mux.HandlePath("GET", "/swagger-ui-bundle.js", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte(staticproto.JS1))
		})
		mux.HandlePath("GET", "/swagger-ui-standalone-preset.js", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte(staticproto.JS2))
		})
		mux.HandlePath("GET", "/swagger-ui.css", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte(staticproto.CSS))
		})
		mux.HandlePath("GET", "/docs.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte(statproto.SwaggerContent))
		})

		log.Printf("server listening at 8081")
		go func() {
			if err := http.ListenAndServe(":8081", mux); err != nil {
				panic(err)
			}
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
