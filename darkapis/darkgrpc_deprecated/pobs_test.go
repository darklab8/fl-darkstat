package darkgrpc_deprecated

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/stretchr/testify/assert"
)

type Api2 struct {
	app_data *appdata.AppData
}

func (a *Api2) GetAppData() *appdata.AppData {
	return a.app_data
}

func RegisterApiRoutes(w *web.Web, app_data *appdata.AppData) *web.Web {
	api := &Api2{
		app_data: app_data,
	}
	api_routes := registry.NewRegister()
	api_routes.Register(GetPobGoodsDeprecated(w, api))

	api_routes.Foreach(func(e *registry.Endpoint) {
		w.GetMux().HandleFunc(string(e.Url), e.Handler)
	})
	return w
}

func TestDepregatedPoBGoods(t *testing.T) {
	ctx := context.Background()

	file_data, err := os.ReadFile(string(utils_os.GetCurrentFolder().Join("testdata", "api_data.json")))
	logus.Log.CheckPanic(err, "failed to read test file")
	desired_result, err := os.ReadFile(string(utils_os.GetCurrentFolder().Join("testdata", "result_desired.json")))
	logus.Log.CheckPanic(err, "failed to read test file")
	new_ctx := context.WithValue(ctx, cfg.CtxKey("pob_goods_data_override"), file_data)
	app_data := router.GetAppDataFixture(new_ctx)
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

	web_closer := web_server.Serve(web.WebServeOpts{SockAddress: some_socket, Port: ptr.Ptr(8433)})

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
					return net.Dial("tcp", "localhost:8433")
				},
			},
		}
	}

	t.Run("GetPoBGoodsDeprecated", func(t *testing.T) {
		post_body, err := json.Marshal(make(map[string]string))
		logus.Log.CheckPanic(err, "unable to marshal post body", typelog.OptError(err))

		res, err := httpc.Post("http://localhost/statproto.Darkstat/GetPoBGoods", ApplicationJson, bytes.NewBuffer(post_body))
		logus.Log.CheckPanic(err, "error making http request: %s\n", typelog.OptError(err))

		resBody, err := io.ReadAll(res.Body)
		logus.Log.CheckPanic(err, "client: could not read response body: %s\n", typelog.OptError(err))

		var result pb.GetPoBGoodsReply
		// fmt.Println("resBody=", string(resBody))

		err = os.WriteFile(string(utils_os.GetCurrentFolder().Join("testdata", "result_received.json")), resBody, os.FileMode(0644))
		logus.Log.CheckPanic(err, "can not write result received", typelog.OptError(err))
		err = json.Unmarshal(resBody, &result)
		logus.Log.CheckPanic(err, "can not unmarshal", typelog.OptError(err))
		assert.Greater(t, len(result.Items), 0)

		// are they equal?
		_ = resBody
		_ = desired_result
		m1 := make(map[string]any)
		m2 := make(map[string]any)
		err1 := json.Unmarshal(resBody, &m1)
		err2 := json.Unmarshal(desired_result, &m2)
		logus.Log.CheckPanic(err1, "failed to unmarshal to m1")
		logus.Log.CheckPanic(err2, "failed to unmarshal to m2")
		// m1_by_key := make(map[string]any)
		// for _, value := range m1["items"] {
		// 	nickname := value["nickname"].(string)
		// 	m1_by_key[nickname] = value

		// }
		// m2_by_key := make(map[string]any)
		// for _, value := range m2.Items {
		// 	nickname := value["nickname"].(string)
		// 	m2_by_key[nickname] = value
		// }
		// for key, value := range m1_by_key {
		// 	value2 := m2_by_key[key]
		// 	are_equal := reflect.DeepEqual(value, value2)
		// 	if !are_equal {
		// 		fmt.Println("not equal!! at nickname=", key)
		// 		fmt.Println(value)
		// 		fmt.Println(value2)
		// 		fmt.Println("not equal END")
		// 	}
		// 	assert.True(t, are_equal)
		// }

		// assert.True(t, reflect.DeepEqual(m1, m2))
		path, v1, v2, found := FirstMismatch(m1, m2, "$")
		if found {
			fmt.Printf("Mismatch at %s: m1=%v, m2=%v\n", path, v1, v2)
			panic("miss match is found")
		} else {
			fmt.Println("Maps are equal")
		}
	})
	web_closer.Close()
}

type GetPoBGoodsReply struct {
	Items []map[string]any `json:"items"`
}

const ApplicationJson = "application/json"

// normalizeSlice sorts slices recursively so order does not matter
func normalizeSlice(slice []any) []any {
	normalized := make([]any, len(slice))
	for i, v := range slice {
		switch vv := v.(type) {
		case []any:
			normalized[i] = normalizeSlice(vv)
		case map[string]any:
			normalized[i] = normalizeMap(vv)
		default:
			normalized[i] = vv
		}
	}
	sort.SliceStable(normalized, func(i, j int) bool {
		return fmt.Sprintf("%#v", normalized[i]) < fmt.Sprintf("%#v", normalized[j])
	})
	return normalized
}

func normalizeMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		switch vv := v.(type) {
		case []any:
			out[k] = normalizeSlice(vv)
		case map[string]any:
			out[k] = normalizeMap(vv)
		default:
			out[k] = vv
		}
	}
	return out
}

// FirstMismatch finds the first primitive mismatch deeply
func FirstMismatch(v1, v2 any, path string) (string, any, any, bool) {
	switch val1 := v1.(type) {

	case map[string]any:
		val2, ok := v2.(map[string]any)
		if !ok {
			// drill down: can't go deeper, return concrete mismatch
			return path, v1, v2, true
		}
		val1 = normalizeMap(val1)
		val2 = normalizeMap(val2)

		// check keys in both maps in a stable way
		keys := make([]string, 0, len(val1)+len(val2))
		seen := map[string]struct{}{}
		for k := range val1 {
			keys = append(keys, k)
			seen[k] = struct{}{}
		}
		for k := range val2 {
			if _, ok := seen[k]; !ok {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		for _, k := range keys {
			subv1, ok1 := val1[k]
			subv2, ok2 := val2[k]
			if !ok1 {
				return FirstMismatch(nil, subv2, fmt.Sprintf("%s.%s", path, k))
			}
			if !ok2 {
				return FirstMismatch(subv1, nil, fmt.Sprintf("%s.%s", path, k))
			}
			if p, sv1, sv2, found := FirstMismatch(subv1, subv2, fmt.Sprintf("%s.%s", path, k)); found {
				return p, sv1, sv2, true
			}
		}
		return "", nil, nil, false

	case []any:
		val2, ok := v2.([]any)
		if !ok {
			return path, v1, v2, true
		}
		val1 = normalizeSlice(val1)
		val2 = normalizeSlice(val2)

		minLen := len(val1)
		if len(val2) < minLen {
			minLen = len(val2)
		}
		for i := 0; i < minLen; i++ {
			if p, sv1, sv2, found := FirstMismatch(val1[i], val2[i], fmt.Sprintf("%s[%d]", path, i)); found {
				return p, sv1, sv2, true
			}
		}
		if len(val1) != len(val2) {
			// mismatch in length -> report as concrete values
			return fmt.Sprintf("%s[len]", path), len(val1), len(val2), true
		}
		return "", nil, nil, false

	default:
		// primitive value check
		if !reflect.DeepEqual(v1, v2) {
			return path, v1, v2, true
		}
		return "", nil, nil, false
	}
}
