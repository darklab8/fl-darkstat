package services

import "github.com/darklab8/fl-darkstat/darkstat/configs_export"

func FilterNicknames[T Nicknamable](filter_nicknames []string, items []T) []T {
	if len(filter_nicknames) == 0 {
		return items
	}

	var result []T
	filter_nicknames_map := make(map[string]bool)
	for _, filter := range filter_nicknames {
		filter_nicknames_map[filter] = true
	}

	for _, item := range items {
		if _, ok := filter_nicknames_map[item.GetNickname()]; ok {
			result = append(result, item)
		}
	}

	return result
}

type Stringable interface {
	comparable
	ToStr() string
}

func FilterMarketGoodCategory[T Stringable](filter_category []string, items map[T]*configs_export.MarketGood) map[string]*configs_export.MarketGood {
	var result map[string]*configs_export.MarketGood = make(map[string]*configs_export.MarketGood)
	if len(filter_category) == 0 {
		for key, item := range items {
			result[key.ToStr()] = item
		}
	}

	filter_category_map := make(map[string]bool)
	for _, filter := range filter_category {
		filter_category_map[filter] = true
	}

	for key, item := range items {
		if _, ok := filter_category_map[item.Category]; ok {
			result[key.ToStr()] = item
		}
	}

	return result
}
