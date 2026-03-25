package app_launch

import (
	"context"

	"github.com/darklab8/fl-darkstat/darkapis/darkrpc_deprecated"
	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

func main() {
	if false {
		// TODO RPC no longer shall be launched.
		// remains purely for as a code documented curiosity how to use it only
		ctx_span := context.Background()
		var rpc_opts []darkrpc_deprecated.ServerOpt
		if settings.Env.EnableUnixSockets {
			rpc_opts = append(rpc_opts, darkrpc_deprecated.WithSockSrv(darkrpc_deprecated.DarkstatRpcSock))
		}
		rpc_server := darkrpc_deprecated.NewRpcServer(rpc_opts...)
		app_data := appdata.NewAppData(ctx_span)
		rpc_server.Serve(app_data)
		if false {
			rpc_server.Close()
		}
	}
}
