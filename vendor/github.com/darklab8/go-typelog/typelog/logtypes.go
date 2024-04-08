package typelog

import (
	"fmt"
	"log/slog"
	"time"
)

type LogAtrs struct {
	slogs []SlogAttr
}

func (s LogAtrs) Render() []SlogAttr {
	return s.slogs
}

func (s *LogAtrs) Append(params ...slog.Attr) {
	for _, param := range params {
		s.slogs = append(s.slogs, param)
	}
}

type LogType func(r *LogAtrs)

type SlogAttr = any

func newSlogArgs(opts ...LogType) []SlogAttr {
	client := &LogAtrs{}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Int("test_param", value))
	}
}

func Any(key string, value any) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.String(key, fmt.Sprintf("%v", value)))
	}
}

func String(key string, value string) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.String(key, value))
	}
}

func Int(key string, value int) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Int(key, value))
	}
}
func Int64(key string, value int64) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Int64(key, value))
	}
}
func Float32(key string, value float32) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Float64(key, float64(value)))
	}
}
func Time(key string, value time.Time) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Time(key, value))
	}
}
func Float64(key string, value float64) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Float64(key, value))
	}
}
func Bool(key string, value bool) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.Bool(key, value))
	}
}

func Expected(value any) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.String("expected", fmt.Sprintf("%v", value)))
	}
}
func Actual(value any) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.String("actual", fmt.Sprintf("%v", value)))
	}
}

func OptError(err error) LogType {
	return func(c *LogAtrs) {
		c.Append(
			slog.String("error_msg", fmt.Sprintf("%v", err)),
			slog.String("error_type", fmt.Sprintf("%T", err)),
		)
	}
}

func Items[T any](item_name string, value []T) LogType {
	return func(c *LogAtrs) {
		sliced_string := fmt.Sprintf("%v", value)
		if len(sliced_string) > 300 {
			sliced_string = sliced_string[:300] + "...sliced string"
		}
		c.Append(slog.String(item_name, fmt.Sprintf("%v", sliced_string)))
		c.Append(slog.String(fmt.Sprintf("%s_len", item_name), fmt.Sprintf("%d", len(value))))
	}
}

func Records[T any](value []T) LogType {
	return Items[T]("records", value)
}

func Args(value []string) LogType {
	return Items[string]("args", value)
}

func Bytes(key string, value []byte) LogType {
	return func(c *LogAtrs) {
		c.Append(slog.String(key, string(value)))
	}
}

type StructType = any

func Struct(value StructType) LogType {
	return func(c *LogAtrs) {
		c.Append(TurnMapToAttrs(StructToMap(value))...)
	}
}

func NestedStruct(key string, value StructType) LogType {
	return func(c *LogAtrs) {
		attrs := TurnMapToAttrs(StructToMap(value))
		c.Append(Group(key, attrs...))
	}
}

func Map(value map[string]any) LogType {
	return func(c *LogAtrs) {
		c.Append(TurnMapToAttrs(value)...)
	}
}

func NestedMap(key string, value map[string]any) LogType {
	return func(c *LogAtrs) {
		c.Append(Group(key, TurnMapToAttrs(value)...))
	}
}
