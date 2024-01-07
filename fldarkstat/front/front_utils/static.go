package front_utils

import (
	"fmt"
	"path/filepath"
)

const (
	StaticRoute string = "/web/static"
)

func GetStatisRoute(path ...string) string {
	fmt.Sprintln("static_route_start")
	result := filepath.Join(append([]string{StaticRoute}, path...)...)
	fmt.Sprintln("static_route=", result)
	return result
}
