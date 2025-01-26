/*
Prototype code, delete once inireader package established itself
*/

package main

import (
	"fmt"
	"strconv"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
)

type INIFile struct {
	File     *file.File
	Sections []Section
}

/*
[BaseGood] // this is Type
abc = 123 // this is Param going into list and hashmap
*/
type Section struct {
	Type   string
	Params []Param
	// denormialization of Param list due to being more comfortable
	ParamMap map[string][]Param
}

// abc = qwe, 1, 2, 3, 4
// abc is key
// qwe is first value
// qwe, 1, 2, 3, 4 are values
// ;abc = qwe, 1, 2, 3 is Comment
type Param struct {
	Key       string
	Values    []UniValue
	IsComment bool     // if commented out
	First     UniValue // denormalization due to very often being needed
}

type UniValue interface {
	AsString() string
}

type ValueString string
type ValueNumber struct {
	value     float64
	precision int
}

func (v ValueString) AsString() string {
	return string(v)
}

func (v ValueNumber) AsString() string {
	return strconv.FormatFloat(float64(v.value), 'f', v.precision, 64)
}

func UniParse(input string) UniValue {

	return ValueString("123")
	// if numberParser.MatchString(input) {

	// }
	// return nil
}

func main() {
	var v UniValue

	v = UniParse("abc")
	printer(v)

	v = UniParse("4")
	printer(v)

	v = UniParse("54.95")
	printer(v)

}

func printer(v UniValue) {
	fmt.Printf("UniValue=(%v, %T)\n", v, v)
	fmt.Printf("AsString=%v\n", v.AsString())

	switch v.(type) {
	case ValueNumber:
		{
			fmt.Println("it is ValueNumber")
		}

	case ValueString:
		{
			fmt.Println("it is ValueString")
		}
	}
}
