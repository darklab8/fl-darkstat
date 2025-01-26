package core_front

import (
	"fmt"
	"testing"
)

func TestCheckFilesystem(t *testing.T) {

	var SomeMap map[string]int = make(map[string]int)

	fmt.Println(SomeMap["not_exisitng_key"])
}
