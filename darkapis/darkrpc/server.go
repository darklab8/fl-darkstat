package darkrpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type RpcServer struct {
	sock_address string
	port         int
}

func NewRpcServer(opts ...ServerOpt) *RpcServer {
	srv := &RpcServer{
		port: 8100,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

type ServerOpt func(r *RpcServer)

func WithSockSrv(sock string) ServerOpt {
	return func(r *RpcServer) {
		r.sock_address = sock
	}
}

func WithPortSrv(port int) ServerOpt {
	return func(r *RpcServer) {
		r.port = port
	}
}

var (
	Log = logus.Log.WithScope("darkgrpc")
)

func (r *RpcServer) Serve(app_data *appdata.AppData) {
	rpcServer := rpc.NewServer()
	rpc_server := NewRpc(app_data)        // NOTE: RPC server is created
	err := rpcServer.Register(rpc_server) // NOTE: RPC server is registered
	Log.CheckError(err, "failed to register rpc server")

	rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath) // NOTE: Handle path
	tcp_address := fmt.Sprintf(":%d", r.port)
	fmt.Println("starting rpc server at ", tcp_address)
	tcp_listener, err := net.Listen("tcp", tcp_address) // NOTE: listen stuff
	if err != nil {
		log.Fatal("listen error:", err)
	}

	settings.Env.EnableUnixSockets = true
	var sock_listener net.Listener
	if cfg.IsLinux && r.sock_address != "" {
		_ = os.Remove(r.sock_address)
		_ = os.Mkdir("/tmp/darkstat", 0777)
		sock_listener, err = net.Listen("unix", r.sock_address) // if serving over Unix
		if err != nil {
			log.Fatal("listen error:", err)
		}
		fmt.Println("turning on server")
		if cfg.IsLinux {
			go rpcServer.Accept(sock_listener) // if serving over Unix

		}
	}

	go func() {
		fmt.Println("http rpc server is launching")
		err := http.Serve(tcp_listener, nil) // NOTE: Server
		if err != nil {
			log.Fatal("http error:", err)
		}
	}()

	fmt.Println("rpc server is launched")
}

func (r *RpcServer) Close() {
	fmt.Println("gracefully existing rpc server")
	logus.Log.CheckError(os.Remove(r.sock_address), "unable removing sock")
}
