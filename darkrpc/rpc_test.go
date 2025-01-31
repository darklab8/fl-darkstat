package darkrpc

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/stretchr/testify/assert"
)

func TestRpc(t *testing.T) {
	some_socket := "/tmp/darkstat/rpc_test.sock"

	app_data := router.NewAppData()
	srv := NewRpcServer(WithSockSrv(some_socket), WithPortSrv(8523))
	srv.Serve(app_data)

	args := Args{}
	client := NewClient(WithSockCli(some_socket))

	t.Run("Get health", func(t *testing.T) {
		reply := new(bool)
		err := client.GetHealth(args, reply)
		logus.Log.CheckPanic(err, "failed to get health")

		assert.NotNil(t, reply)
		assert.True(t, *reply)
	})

	// Setup code for given condition goes here
	t.Run("Check API works", func(t *testing.T) {
		var reply Reply
		err := client.GetBases(args, &reply)
		logus.Log.CheckPanic(err, "failed to get bases")
		fmt.Println("Bases[0]=", reply.Bases[0])
	})

	// Teardown code for given condition goes here
	srv.Close()
}
