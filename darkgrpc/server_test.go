package darkgrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRpc(t *testing.T) {

	test_port := 8524

	app_data := router.GetAppDataFixture()

	// TODO write some day sock support
	// some_socket := "/tmp/darkstat/api_test.sock"

	grpc_server := NewServer(app_data, test_port)
	go grpc_server.Serve()

	c := NewClient(fmt.Sprintf("localhost:%d", test_port))
	defer c.Conn.Close()
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 10e6)

	t.Run("GetHealth", func(t *testing.T) {
		res, err := c.GetHealth(context.Background(), &statproto.Empty{})
		logus.Log.CheckPanic(err, "error making rpc call: %s\n", typelog.OptError(err))

		assert.True(t, res.IsHealthy)
	})

	t.Run("GetBases", func(t *testing.T) {
		res, err := c.GetBases(context.Background(), &statproto.GetBasesInput{IncludeMarketGoods: true}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get bases: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetCommodities", func(t *testing.T) {
		res, err := c.GetCommodities(context.Background(), &statproto.GetCommoditiesInput{IncludeMarketGoods: true}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get commoditieis: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetAmmos", func(t *testing.T) {
		res, err := c.GetAmmos(context.Background(), &statproto.GetEquipmentInput{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get ammos: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetCounterMeasures", func(t *testing.T) {
		res, err := c.GetCounterMeasures(context.Background(), &statproto.GetEquipmentInput{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get cms: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetEngines", func(t *testing.T) {
		res, err := c.GetEngines(context.Background(), &statproto.GetEquipmentInput{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get cms: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})
}
