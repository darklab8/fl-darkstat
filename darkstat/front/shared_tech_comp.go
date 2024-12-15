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

func GetTdDiscoCacheKey(disco types.DiscoveryIDs, nickname string) TdCacheKey {
	cache_key_data := marshalIDs(disco, nickname)
	h := md5.New()
	io.WriteString(h, cache_key_data)
	return TdCacheKey(fmt.Sprintf("%x", h.Sum(nil)))
}

func GetDiscoCacheMap(items []Item, disco types.DiscoveryIDs) map[TdCacheKey]DiscoCompat {
	var cache map[TdCacheKey]DiscoCompat = map[TdCacheKey]DiscoCompat{}
	for _, item := range items {
		nickname := item.GetNickname()
		cache_key := GetTdDiscoCacheKey(disco, nickname)
		cache[cache_key] = DiscoCompat{
			Nickname: nickname,
			Data:     item.GetTechCompat(),
		}
	}
	return cache
}
