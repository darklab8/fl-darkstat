package darkhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/stretchr/testify/assert"
)

type TestOpts struct {
	CheckMarketGoods bool
	CheckTechCompat  bool
	IncludeRephacks  bool
}

func FixtureTestItems[T Nicknamable](t *testing.T, httpc http.Client, url string, test_name string, opts TestOpts) []T {
	body := map[string]any{}

	if opts.CheckMarketGoods {
		body["include_market_goods"] = true
	}
	if opts.CheckTechCompat {
		body["include_tech_compat"] = true
	}
	if opts.IncludeRephacks {
		body["include_rephacks"] = true
	}
	post_body, err := json.Marshal(body)
	logus.Log.CheckError(err, "failing to marshal body")
	res, err := httpc.Post("http://localhost/api"+url, ApplicationJson, bytes.NewBuffer(post_body))
	logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

	var items []T
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		fmt.Println("resBody=", string(resBody))
	}
	logus.Log.CheckPanic(err, "can not unmarshal", typelog.OptError(err))

	assert.Greater(t, len(items), 0)
	fmt.Println(items[0])

	return items
}

func TestApiHealth(t *testing.T) {

	app_data := &appdata.AppData{
		Configs: &configs_export.Exporter{},
	}
	stat_fs := &builder.Filesystem{}

	some_socket := "/tmp/darkstat/api_test2.sock"
	server_port := 8454

	web_server := RegisterApiRoutes(web.NewWeb(
		[]*builder.Filesystem{
			stat_fs,
		},
		web.WithSiteRoot(settings.Env.SiteRoot),
	), app_data)

	var web_closer web.ServerClose
	var httpc http.Client
	if settings.Env.EnableUnixSockets {
		web_closer = web_server.Serve(web.WebServeOpts{SockAddress: some_socket})
		httpc = http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", some_socket)
				},
			},
		}
	} else {
		web_closer = web_server.Serve(web.WebServeOpts{Port: ptr.Ptr(server_port)})
		httpc = http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("tcp", fmt.Sprintf("localhost:%d", server_port))
				},
			},
		}
	}

	t.Run("GetHealth", func(t *testing.T) {
		res, err := httpc.Get("http://localhost/ping")
		logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

		resBody, err := io.ReadAll(res.Body)
		logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

		assert.Contains(t, string(resBody), "pong")

		t.Run("GetHealthSubtest", func(t *testing.T) {
			assert.True(t, true)
			fmt.Println("executed subtest too")
		})
	})
	web_closer.Close()

}

func TestApi(t *testing.T) {
	ctx := context.Background()
	app_data := router.GetAppDataFixture()
	stat_router := router.NewRouter(app_data)
	stat_builder := stat_router.Link(ctx)
	stat_fs := stat_builder.BuildAll(true, nil)

	var some_socket string
	if settings.Env.EnableUnixSockets {
		some_socket = "/tmp/darkstat/api_test.sock"
	}

	web_server := RegisterApiRoutes(web.NewWeb(
		[]*builder.Filesystem{
			stat_fs,
		},
		web.WithMutexableData(app_data),
		web.WithSiteRoot(settings.Env.SiteRoot),
	), app_data)

	web_closer := web_server.Serve(web.WebServeOpts{SockAddress: some_socket, Port: ptr.Ptr(8432)})

	var httpc http.Client
	if settings.Env.EnableUnixSockets {
		httpc = http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", some_socket)
				},
			},
		}
	} else {
		httpc = http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("tcp", "localhost:8432")
				},
			},
		}
	}

	t.Run("GetHealth", func(t *testing.T) {
		res, err := httpc.Get("http://localhost/ping")
		logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

		resBody, err := io.ReadAll(res.Body)
		logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

		assert.Contains(t, string(resBody), "pong")
	})

	t.Run("GetBases", func(t *testing.T) {
		items := FixtureTestItems[Base](t, httpc, "/npc_bases", "NpcBases", TestOpts{
			CheckMarketGoods: true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)

		t.Run("GetGraphPaths", func(t *testing.T) {
			nicknames := []appdata.GraphPathReq{
				{
					From: string(items[0].Nickname),
					To:   string(items[1].Nickname),
				},
			}

			post_body, err := json.Marshal(nicknames)
			logus.Log.CheckPanic(err, "unable to marshal post body", typelog.OptError(err))

			res, err := httpc.Post("http://localhost/api/graph/paths", ApplicationJson, bytes.NewBuffer(post_body))
			logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

			resBody, err := io.ReadAll(res.Body)
			logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

			var items []appdata.GraphPathsResp
			fmt.Println("resBody=", string(resBody))
			err = json.Unmarshal(resBody, &items)
			logus.Log.CheckPanic(err, "can not unmarshal", typelog.OptError(err))

			assert.Greater(t, len(items), 0)
			assert.Equal(t, 1, len(items))

			if items[0].Error != nil {
				fmt.Println("items[0].Error=", *items[0].Error)
			}
			assert.Nil(t, items[0].Error)
		})

		t.Run("GetInfocards", func(t *testing.T) {
			var nickname []string = []string{
				items[0].GetNickname(),
				items[1].GetNickname(),
				"smth that does not exist",
			}

			post_body, err := json.Marshal(nickname)
			logus.Log.CheckPanic(err, "unable to marshal post body", typelog.OptError(err))

			res, err := httpc.Post("http://localhost/api"+"/infocards", ApplicationJson, bytes.NewBuffer(post_body))
			logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

			resBody, err := io.ReadAll(res.Body)
			logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

			var items []InfocardResp
			fmt.Println("resBody=", string(resBody))
			err = json.Unmarshal(resBody, &items)
			logus.Log.CheckPanic(err, "can not unmarshal", typelog.OptError(err))

			assert.Greater(t, len(items), 0)

			assert.Nil(t, items[0].Error)
			assert.Nil(t, items[1].Error)
			assert.NotNil(t, items[2].Error)
		})
	})

	t.Run("GetCommodities", func(t *testing.T) {
		items := FixtureTestItems[Commodity](t, httpc, "/commodities", "Commodities", TestOpts{
			CheckMarketGoods: true,
		})
		has_market_goods := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				has_market_goods = true
			}
		}
		assert.True(t, has_market_goods)
	})

	t.Run("GetFactions", func(t *testing.T) {
		_ = FixtureTestItems[configs_export.Faction](t, httpc, "/factions", "Factions", TestOpts{})
	})

	t.Run("GetPoBs", func(t *testing.T) {
		if !app_data.Configs.IsDiscovery {
			return
		}
		_ = FixtureTestItems[configs_export.PoB](t, httpc, "/pobs", "Pobs", TestOpts{})
		_ = FixtureTestItems[configs_export.Base](t, httpc, "/pobs/bases", "PoBBases", TestOpts{})
	})

	t.Run("GetPoBGoods", func(t *testing.T) {
		if !app_data.Configs.IsDiscovery {
			return
		}
		_ = FixtureTestItems[configs_export.PoBGood](t, httpc, "/pob_goods", "PoBGodds", TestOpts{})
	})

	t.Run("GetShips", func(t *testing.T) {
		items := FixtureTestItems[Ship](t, httpc, "/ships", "Ships", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		has_market_goods := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				has_market_goods = true
			}
		}
		assert.True(t, has_market_goods)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetTractors", func(t *testing.T) {
		if !app_data.Configs.IsDiscovery {
			return
		}
		items := FixtureTestItems[Tractor](t, httpc, "/tractors", "Tractors", TestOpts{
			CheckMarketGoods: true,
			IncludeRephacks:  true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)

		var found Tractor
		for _, tractor := range items {
			if len(tractor.Rephacks) > 0 {
				found = tractor
				break
			}
		}
		assert.Greater(t, len(found.Rephacks), 0)
	})

	t.Run("GetAmmos", func(t *testing.T) {
		items := FixtureTestItems[Ammo](t, httpc, "/ammos", "Ammos", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetGuns", func(t *testing.T) {
		for _, gun := range app_data.Configs.Guns {
			_, err := json.Marshal(gun)
			assert.Nil(t, err)
		}

		items := FixtureTestItems[Gun](t, httpc, "/guns", "Ships", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		has_market_goods := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				has_market_goods = true
			}
		}

		assert.True(t, has_market_goods)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetMissiles", func(t *testing.T) {
		items := FixtureTestItems[Gun](t, httpc, "/missiles", "Missiles", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})

		market_goods_exist := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				market_goods_exist = true
			}
		}
		if !market_goods_exist {
			panic("not a single market good")
		}
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetMines", func(t *testing.T) {
		items := FixtureTestItems[Mine](t, httpc, "/mines", "Mines", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		market_goods_exist := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				market_goods_exist = true
			}
		}
		if !market_goods_exist {
			panic("not a single market good")
		}
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetCMs", func(t *testing.T) {
		items := FixtureTestItems[CounterMeasure](t, httpc, "/counter_measures", "CounterMeasures", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetEngines", func(t *testing.T) {
		if !app_data.Configs.IsDiscovery {
			return
		}
		for _, item := range app_data.Configs.Engines {
			_, err := json.Marshal(item)
			if err != nil {
				fmt.Println("found broken engine")
			}
			assert.Nil(t, err)
		}
		items := FixtureTestItems[Engine](t, httpc, "/engines", "Engines", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		has_market_goods := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				has_market_goods = true
			}
		}
		assert.True(t, has_market_goods)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetScanners", func(t *testing.T) {
		if !app_data.Configs.IsDiscovery {
			return
		}
		items := FixtureTestItems[Scanner](t, httpc, "/scanners", "Scanners", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		has_market_goods := false
		for _, item := range items {
			if len(item.MarketGoods) > 0 {
				has_market_goods = true
			}
		}
		assert.True(t, has_market_goods)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetShields", func(t *testing.T) {
		items := FixtureTestItems[Shield](t, httpc, "/shields", "Shields", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	t.Run("GetThrusters", func(t *testing.T) {
		for _, item := range app_data.Configs.Thrusters {
			_, err := json.Marshal(item)
			if err != nil {
				fmt.Println("found broken engine")
			}
			assert.Nil(t, err)
		}
		items := FixtureTestItems[Thruster](t, httpc, "/thrusters", "Thrusters", TestOpts{
			CheckMarketGoods: true,
			CheckTechCompat:  true,
		})
		assert.Greater(t, len(items[0].MarketGoods), 0)
		if app_data.Configs.Mapped.Discovery != nil {
			has_tech_compat := false
			for _, item := range items {
				if len(item.TechCompat.TechcompatByID) > 0 {
					has_tech_compat = true
				}
			}
			assert.True(t, has_tech_compat)
		}
	})

	// // Teardown code for given condition goes here
	web_closer.Close()
}

const ApplicationJson = "application/json"
