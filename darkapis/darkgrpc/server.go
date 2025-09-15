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
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	_ "github.com/darklab8/fl-darkstat/docs"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
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
	reg          *prometheus.Registry
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
func WithProm(reg *prometheus.Registry) ServerOpt {
	return func(s *Server) {
		s.reg = reg
	}
}

type ServerOpt func(s *Server)

const DefaultServerPort = 50051

const DarkstatGRpcSock = "/tmp/darkstat/grpc.sock"

func (r *Server) Serve() {
	// Setup metrics.
	// inspired by https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/providers/prometheus
	// and https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/examples/server/main.go
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(prometheus.DefBuckets),
		),
	)
	if r.reg != nil {
		r.reg.MustRegister(srvMetrics)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(32*10e6),
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(),
		),
	)
	pb.RegisterDarkstatServer(s, r)
	log.Printf("server listening at %v", lis.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(s)

	srvMetrics.InitializeMetrics(s)

	if cfg.IsLinux && settings.Env.EnableUnixSockets && r.sock_address != "" {
		_ = os.Remove(r.sock_address)
		_ = os.Mkdir("/tmp/darkstat", 0777)
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
		if cfg.IsLinux && settings.Env.EnableUnixSockets {
			err = pb.RegisterDarkstatHandlerFromEndpoint(ctx, mux, fmt.Sprintf("unix:%s", DarkstatGRpcSock), opts)
		} else {
			err = pb.RegisterDarkstatHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
		}
		if err != nil {

			panic(err)
		}

		// Sprinked with API documentation :)
		err = mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, err = w.Write([]byte(staticproto.Index))
			logus.Log.CheckError(err, "failed to write index file")
		})
		logus.Log.CheckError(err, "mux failed to handle path for index file")
		err = mux.HandlePath("GET", "/swagger-ui-bundle.js", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, err = w.Write([]byte(staticproto.JS1))
			logus.Log.CheckError(err, "failed to write  file JS1 for swagger")
		})
		logus.Log.CheckError(err, "mux failed to handle path for JS1 file")

		err = mux.HandlePath("GET", "/swagger-ui-standalone-preset.js", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, err = w.Write([]byte(staticproto.JS2))
			logus.Log.CheckError(err, "failed to write  file JS2 for swagger")
		})
		logus.Log.CheckError(err, "mux failed to handle path for JS2 file")

		err = mux.HandlePath("GET", "/swagger-ui.css", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, err = w.Write([]byte(staticproto.CSS))
			logus.Log.CheckError(err, "failed to write  file CSS for swagger")
		})
		logus.Log.CheckError(err, "mux failed to handle path for CSS file for swagger")

		err = mux.HandlePath("GET", "/docs.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			_, err = w.Write([]byte(statproto.SwaggerContent))
			logus.Log.CheckError(err, "failed to write  file docs.json for swagger")
		})
		logus.Log.CheckError(err, "mux failed to handle path for docs.json for swagger")

		log.Printf("server listening at 8081")
		s := &http.Server{
			Addr:    ":8081",
			Handler: MiddlewarePrometheusForAPIGateway(mux),
		}

		go func() {
			if err := s.ListenAndServe(); err != nil {
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
		s.app_data.RLock()
		defer s.app_data.RUnlock()
	}

	return &pb.HealthReply{
		IsHealthy: true,
	}, nil
}
