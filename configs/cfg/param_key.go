package cfg

type ParamKey string

type Key = ParamKey

// var contexts = sync.Map{}

// func Key(key string) ParamKey {
// 	value, ok := contexts.Load(key)
// 	if !ok {
// 		value = ParamKey(ptr.Ptr(key))
// 		contexts.Store(key, value)
// 	}
// 	return value.(ParamKey)
// }
