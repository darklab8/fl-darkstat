package front

import (
	"crypto/md5"
	"fmt"
	"io"

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

func GetTdDiscoCacheKey(shared types.SharedData, nickname string) TdCacheKey {
	if !shared.ShowDisco {
		return ""
	}
	cache_key_data := marshalIDs(shared, nickname)
	h := md5.New()
	io.WriteString(h, cache_key_data)
	return TdCacheKey(fmt.Sprintf("%x", h.Sum(nil)))
}

func GetDiscoCacheMap(items []Item, shared types.SharedData) map[TdCacheKey]DiscoCompat {
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
