package services

// import "github.com/darklab8/fl-darkstat/darkstat/configs_export"

// func FilterNicknames[T Nicknamable](filter_nicknames []string, items []T) []T {
// 	if len(filter_nicknames) == 0 {
// 		return items
// 	}

// 	var result []T
// 	filter_nicknames_map := make(map[string]bool)
// 	for _, filter := range filter_nicknames {
// 		filter_nicknames_map[filter] = true
// 	}

// 	for _, item := range items {
// 		if _, ok := filter_nicknames_map[item.GetNickname()]; ok {
// 			result = append(result, item)
// 		}
// 	}

// 	return result
// }

// func FilterMarketGoodCategory[T comparable](filter_category []string, items map[T]*configs_export.MarketGood) []*configs_export.MarketGood {
// 	var result []*configs_export.MarketGood
// 	if len(filter_category) == 0 {
// 		for _, item := range items {
// 			result = append(result, item)
// 		}
// 	}

// 	filter_category_map := make(map[string]bool)
// 	for _, filter := range filter_category {
// 		filter_category_map[filter] = true
// 	}

// 	for _, item := range items {
// 		if _, ok := filter_category_map[item.Category]; ok {
// 			result = append(result, item)
// 		}
// 	}

// 	return result
// }
