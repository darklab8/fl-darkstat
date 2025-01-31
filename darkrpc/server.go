package darkrpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type RpcServer struct {
}

func NewRpcServer() RpcServer { return RpcServer{} }

func (r *RpcServer) Serve(app_data *router.AppData) {
	rpcServer := rpc.NewServer()

	rpc_server := NewRpc(app_data)
	rpcServer.Register(rpc_server)

	rpc.HandleHTTP()                                // if serving over http
	tcp_listener, err := net.Listen("tcp", ":8100") // if serving over http
	if err != nil {
		log.Fatal("listen error:", err)
	}

	var sock_listener net.Listener
	if cfg.IsLinux {
		os.Remove(DarkstatSock)
		os.Mkdir("/tmp/darkstat", 0777)
		sock_listener, err = net.Listen("unix", DarkstatSock) // if serving over Unix
		if err != nil {
			log.Fatal("listen error:", err)
		}
		fmt.Println("turning on server")
		if cfg.IsLinux {
			go rpcServer.Accept(sock_listener) // if serving over Unix

		}
	}

	go func() {
		err := http.Serve(tcp_listener, nil) // if serving over Http
		if err != nil {
			log.Fatal("http error:", err)
		}
	}()

	fmt.Println("rpc server is launched")
}

func (r *RpcServer) Close() {
	fmt.Println("gracefully existing rpc server")
	logus.Log.CheckError(os.Remove(DarkstatSock), "unable removing sock")
}
