package disco

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
)

type TdCacheKey string

type DiscoCompat struct {
	Nickname string
	Data     *configs_export.DiscoveryTechCompat
}

type Item interface {
	GetNickname() string
	GetTechCompat() *configs_export.DiscoveryTechCompat
}

func GetTdDiscoCacheKey(shared *types.SharedData, nickname string) TdCacheKey {
	if !shared.ShowDisco {
		return ""
	}
	cache_key_data := marshalIDs(shared, nickname)
	h := md5.New()
	io.WriteString(h, cache_key_data)
	return TdCacheKey(fmt.Sprintf("%x", h.Sum(nil)))
}

func GetDiscoCacheMap(items []Item, shared *types.SharedData) map[TdCacheKey]DiscoCompat {
	var cache map[TdCacheKey]DiscoCompat = map[TdCacheKey]DiscoCompat{}

	if !shared.ShowDisco {
		return cache
	}
	for _, item := range items {
		nickname := item.GetNickname()
		cache_key := GetTdDiscoCacheKey(shared, nickname)
		cache[cache_key] = DiscoCompat{
			Nickname: nickname,
			Data:     item.GetTechCompat(),
		}
	}
	return cache
}

func marshalIDs(shared *types.SharedData, item_nickname string) string {

	var compat_by_id map[string]float64 = make(map[string]float64)

	compat_by_id[""] = shared.Config.GetCompatibilty(item_nickname, "")

	for _, id := range shared.Ids {
		compat := shared.Config.GetCompatibilty(item_nickname, id.Nickname)

		// data size saving
		if compat <= 0.1 {
			continue
		}

		compat_by_id[string(id.ShortNickname)] = compat
	}

	bytes, _ := json.Marshal(compat_by_id)
	return strings.ReplaceAll(string(bytes), "\"", "'")
}
