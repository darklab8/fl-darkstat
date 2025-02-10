package darkgrpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/stretchr/testify/assert"
)

func TestRpc(t *testing.T) {

	test_port := 8524
	unix_socket := "/tmp/testing_grpc.sock"

	app_data := router.GetAppDataFixture()

	// TODO write some day sock support
	// some_socket := "/tmp/darkstat/api_test.sock"

	grpc_server := NewServer(app_data, test_port, WithSockAddr(unix_socket))
	go grpc_server.Serve()

	// c := NewClient(fmt.Sprintf("unix:%s", unix_socket))
	c := NewClient(fmt.Sprintf("localhost:%d", test_port))
	defer c.Conn.Close()

	for i := 0; i < 10; i++ {
		res, _ := c.GetHealth(context.Background(), &statproto.Empty{})
		if res != nil {
			if res.IsHealthy {
				break
			}
		}
		fmt.Println("test server is not ready yet. Sleeping")
		time.Sleep(5 * time.Second)
	}

	t.Run("GetHealth", func(t *testing.T) {
		res, err := c.GetHealth(context.Background(), &statproto.Empty{})
		logus.Log.CheckPanic(err, "error making rpc call: %s\n", typelog.OptError(err))

		assert.True(t, res.IsHealthy)
	})

	t.Run("GetBases", func(t *testing.T) {
		res, err := c.GetBasesNpc(context.Background(), &statproto.GetBasesInput{IncludeMarketGoods: true})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)

		t.Run("GetGraphPaths", func(t *testing.T) {
			res, err := c.GetGraphPaths(context.Background(), &statproto.GetGraphPathsInput{
				Queries: []*statproto.GraphPathQuery{{
					From: string(res.Items[0].Nickname),
					To:   string(res.Items[1].Nickname),
				}},
			})
			logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
			assert.Greater(t, len(res.Answers), 0)
			assert.Equal(t, 1, len(res.Answers))
			assert.Nil(t, res.Answers[0].Error)
			assert.Greater(t, *res.Answers[0].Time.Transport, int64(0))
		})
	})

	t.Run("GetBasesMiningOperations", func(t *testing.T) {
		res, err := c.GetBasesMiningOperations(context.Background(), &statproto.GetBasesInput{IncludeMarketGoods: true})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetCommodities", func(t *testing.T) {
		res, err := c.GetCommodities(context.Background(), &statproto.GetCommoditiesInput{IncludeMarketGoods: true})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetAmmos", func(t *testing.T) {
		res, err := c.GetAmmos(context.Background(), &statproto.GetEquipmentInput{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetCounterMeasures", func(t *testing.T) {
		res, err := c.GetCounterMeasures(context.Background(), &statproto.GetEquipmentInput{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetEngines", func(t *testing.T) {
		res, err := c.GetEngines(context.Background(), &statproto.GetEquipmentInput{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})

	t.Run("GetPoBs", func(t *testing.T) {
		res, err := c.GetPoBs(context.Background(), &statproto.Empty{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		if app_data.Configs.Configs.Discovery != nil {
			assert.Greater(t, len(res.Items), 0)
		}
	})
	t.Run("GetPoBGoods", func(t *testing.T) {
		res, err := c.GetPoBGoods(context.Background(), &statproto.Empty{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		if app_data.Configs.Configs.Discovery != nil {
			assert.Greater(t, len(res.Items), 0)
		}
	})

	t.Run("GetPoBBases", func(t *testing.T) {
		res, err := c.GetBasesPoBs(context.Background(), &statproto.GetBasesInput{IncludeMarketGoods: true})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		if app_data.Configs.Configs.Discovery != nil {
			assert.Greater(t, len(res.Items), 0)
		}
	})

	t.Run("GetHashes", func(t *testing.T) {
		res, err := c.GetHashes(context.Background(), &statproto.Empty{})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.HashesByNick), 0)
	})

	t.Run("GetFactions", func(t *testing.T) {
		res, err := c.GetFactions(context.Background(), &statproto.GetFactionsInput{
			IncludeReputations: true,
			IncludeBribes:      true,
		})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})
	t.Run("GetThrusters", func(t *testing.T) {
		res, err := c.GetThrusters(context.Background(), &statproto.GetEquipmentInput{
			IncludeMarketGoods: true,
			IncludeTechCompat:  true,
		})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})
	t.Run("GetShips", func(t *testing.T) {
		res, err := c.GetShips(context.Background(), &statproto.GetEquipmentInput{
			IncludeMarketGoods: true,
			IncludeTechCompat:  true,
		})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})
	t.Run("GetShields", func(t *testing.T) {
		res, err := c.GetShields(context.Background(), &statproto.GetEquipmentInput{
			IncludeMarketGoods: true,
			IncludeTechCompat:  true,
		})
		logus.Log.CheckPanic(err, "error making rpc call to get items: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
	})
}
