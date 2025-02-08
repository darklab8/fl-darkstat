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

func FixtureMarketGoodsTest(
	t *testing.T, c *Client,
	method func(ctx context.Context, in *statproto.GetMarketGoodsInput, opts ...grpc.CallOption) (*statproto.GetMarketGoodsReply, error),
	test_name string,
	item1 Nicknamable, item2 Nicknamable) {
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 10e6)

	t.Run("Get"+test_name+"MarketGoods", func(t *testing.T) {
		var nickname []string = []string{
			item1.GetNickname(),
			item2.GetNickname(),
		}

		res, err := method(context.Background(), &statproto.GetMarketGoodsInput{
			Nicknames: nickname,
		}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get market goods: %s\n", typelog.OptError(err))

		answers := res.Answers

		assert.Greater(t, len(answers), 0)

		assert.Nil(t, answers[0].Error)
		assert.Nil(t, answers[1].Error)
	})
}

func FixtureTechCompatTest(
	t *testing.T, c *Client,
	method func(ctx context.Context, in *statproto.GetTechCompatInput, opts ...grpc.CallOption) (*statproto.GetTechCompatReply, error),
	test_name string,
	item1 Nicknamable, item2 Nicknamable) {
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 10e6)

	t.Run("Get"+test_name+"MarketGoods", func(t *testing.T) {
		var nickname []string = []string{
			item1.GetNickname(),
			item2.GetNickname(),
		}

		res, err := method(context.Background(), &statproto.GetTechCompatInput{
			Nicknames: nickname,
		}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get tech compat: %s\n", typelog.OptError(err))

		answers := res.Answers

		assert.Greater(t, len(answers), 0)

		assert.Nil(t, answers[0].Error)
		assert.Nil(t, answers[1].Error)
	})
}

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
		res, err := c.GetBases(context.Background(), &statproto.Empty{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get bases: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
		FixtureMarketGoodsTest(t, c, c.GetBasesMarketGoods, "Bases", res.Items[0], res.Items[1])
	})

	t.Run("GetCommodities", func(t *testing.T) {
		res, err := c.GetCommodities(context.Background(), &statproto.Empty{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get commoditieis: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
		FixtureMarketGoodsTest(t, c, c.GetCommoditiesMarketGoods, "Commodities", res.Items[0], res.Items[1])
	})

	t.Run("GetAmmos", func(t *testing.T) {
		res, err := c.GetAmmos(context.Background(), &statproto.Empty{}, maxSizeOption)
		logus.Log.CheckPanic(err, "error making rpc call to get commoditieis: %s\n", typelog.OptError(err))
		assert.Greater(t, len(res.Items), 0)
		FixtureMarketGoodsTest(t, c, c.GetAmmosMarketGoods, "Ammos", res.Items[0], res.Items[1])
		FixtureTechCompatTest(t, c, c.GetAmmosTechCompat, "Ammos", res.Items[0], res.Items[1])
	})
}
