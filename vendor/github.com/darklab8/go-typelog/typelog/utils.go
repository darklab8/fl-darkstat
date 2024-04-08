package typelog

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

func CompL[T any, V any](objs []T, lambda func(x T) V) []V {
	results := make([]V, 0, len(objs))
	for _, obj := range objs {
		results = append(results, lambda(obj))
	}
	return results
}

func logGroupFiles() slog.Attr {
	return slog.Group("files",
		"file3", GetCallingFile(3),
		"file4", GetCallingFile(4),
	)
}

func GetCallingFile(level int) string {
	GetTwiceParentFunctionLocation := level
	_, filename, _, _ := runtime.Caller(GetTwiceParentFunctionLocation)
	filename = filepath.Base(filename)
	return fmt.Sprintf("f:%s ", filename)
}

func StructToMap(somestruct any) map[string]any {
	var mapresult map[string]interface{}
	inrec, _ := json.Marshal(somestruct)
	json.Unmarshal(inrec, &mapresult)
	return mapresult
}

func TurnMapToAttrs(params map[string]any) []slog.Attr {
	anies := []slog.Attr{}
	for key, value := range params {
		switch v := value.(type) {
		case string:
			anies = append(anies, slog.String(key, v))
		case int:
			anies = append(anies, slog.Int(key, v))
		case int64:
			anies = append(anies, slog.Int64(key, v))
		case float64:
			anies = append(anies, slog.Float64(key, v))
		case float32:
			anies = append(anies, slog.Float64(key, float64(v)))
		case bool:
			anies = append(anies, slog.Bool(key, v))
		case time.Time:
			anies = append(anies, slog.Time(key, v))
		case map[string]any:
			anies = append(anies, Group(key, TurnMapToAttrs(v)...))
		default:
			anies = append(anies, slog.String(key, fmt.Sprintf("%v", v)))
		}
	}

	return anies
}

func TurnStructToAttrs(somestruct any) []slog.Attr {
	return TurnMapToAttrs(StructToMap(somestruct))
}

func Group(name string, attrs ...slog.Attr) slog.Attr {
	return slog.Group(name, CompL(attrs, func(x slog.Attr) SlogAttr { return SlogAttr(x) })...)
}
