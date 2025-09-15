package darkrpc

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/stretchr/testify/assert"
)

func TestRpc(t *testing.T) {
	app_data := router.GetAppDataFixture()

	var srv *RpcServer
	var client RpcI
	if settings.Env.EnableUnixSockets {
		some_socket := "/tmp/darkstat/rpc_test.sock"
		srv = NewRpcServer(WithSockSrv(some_socket), WithPortSrv(8523))
		srv.Serve(app_data)
		client = NewClient(WithSockCli(some_socket))
	} else {
		srv = NewRpcServer(WithPortSrv(8523))
		srv.Serve(app_data)
		client = NewClient(WithPortCli(8523))
	}

	args := Args{}
	t.Run("GetHealth", func(t *testing.T) {
		reply := new(bool)
		err := client.GetHealth(args, reply)
		logus.Log.CheckPanic(err, "failed to get health")

		assert.NotNil(t, reply)
		assert.True(t, *reply)
	})

	// Setup code for given condition goes here
	t.Run("GetBaseCheck", func(t *testing.T) {
		var reply Reply
		err := client.GetBases(args, &reply)
		logus.Log.CheckPanic(err, "failed to get bases")
		fmt.Println("Bases[0]=", reply.Bases[0])
	})

	t.Run("GetInfo", func(t *testing.T) {
		var reply GetInfoReply
		err := client.GetInfo(GetInfoArgs{Query: "Akabat"}, &reply)
		logus.Log.CheckPanic(err, "failed to get info")
		fmt.Println("Content=", reply.Content)
	})

	// Teardown code for given condition goes here
	srv.Close()
}
