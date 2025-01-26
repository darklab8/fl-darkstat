package bini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	gbp "github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped/go-binary-pack"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/bin"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
	"golang.org/x/text/encoding/charmap"
)

type EntryName string
type EntryValues = []interface{}

type Row map[EntryName]EntryValues

type Section struct {
	section_name string
	rows         []Row
}

var bp = new(gbp.BinaryPack)

const SEEK_SET = io.SeekStart // python default seek(offset, whence=os.SEEK_SET, /)

var VALUE_TYPES map[int]string = map[int]string{
	1: "i", 2: "f", 3: "i",
}

// maps a byte value type to a struct format string

func parse_file(path utils_types.FilePath, FoldValues FoldValues) []Section {
	mem := bin.NewBDatas()
	var result []Section = make([]Section, 0, 100)

	var string_table map[int]string = make(map[int]string)
	data, err := os.ReadFile(path.ToString())
	logus.Log.CheckPanic(err, "uunable to open file", utils_logus.FilePath(path))

	file_size := len(data)

	fh := bytes.NewReader(data)

	bdata := make([]byte, 12)
	fh.Read(bdata)

	format := []string{"4s", "I", "I"}

	packed_values, err := bp.UnPack(format, bdata)
	magic := packed_values[0].(string)
	version := packed_values[1].(int)
	str_table_offset := packed_values[2].(int)

	if magic != "BINI" || version != 1 {
		logus.Log.Panic("Expected finding BINI. Found smth else", typelog.String("magic", magic), typelog.Int("version", version))
	}

	fh.Seek(int64(str_table_offset), SEEK_SET)

	var raw_table []byte
	raw_table_length := file_size - str_table_offset - 1
	if raw_table_length <= 0 {
		return result
	}
	raw_table = make([]byte, raw_table_length)
	fh.Read(raw_table)

	raw_tables := bytes.Split(raw_table, []byte{'\x00'})

	count := 0
	for _, table := range raw_tables {

		tr := charmap.Windows1252.NewDecoder().Reader(strings.NewReader(string(table)))
		windows_decoded, err := io.ReadAll(tr)
		logus.Log.CheckPanic(err, "failed decoding to 1252", utils_logus.FilePath(path))

		string_table[count] = string(windows_decoded) // to lower
		count += len(table) + 1
	}

	// return to end of header to read sections
	var position int
	pos, err := fh.Seek(12, SEEK_SET)
	position = int(pos)

	for position < str_table_offset {
		packed_values, offset := bin.Read2(fh, mem.GetBData(4), []string{"h", "h"})
		position += offset
		section_name_ptr := packed_values[0].(int)
		entry_count := packed_values[1].(int)

		section_name := string_table[section_name_ptr]

		var section []Row
		for e := 0; e < entry_count; e++ {
			packed_values, offset := bin.Read2(fh, mem.GetBData(3), []string{"h", "b"})
			position += offset
			entry_name_ptr := packed_values[0].(int)
			value_count := packed_values[1].(int)
			entry_name := string_table[entry_name_ptr]

			var row Row = make(Row)
			row[EntryName(entry_name)] = make([]interface{}, 0, 10)

			for v := 0; v < value_count; v++ {
				packed_values, offset := bin.Read2(fh, mem.GetBData(1), []string{"b"})
				position += offset
				value_type := packed_values[0].(int)

				packed_values, offset = bin.Read2(fh, mem.GetBData(4), []string{VALUE_TYPES[value_type]})
				position += offset

				var value_data interface{}
				switch len(packed_values) {
				case 1:
					value_data = packed_values[0]
				case 0:
					//pass
				default:
					logus.Log.Panic("expected 1 or 0 packed values", typelog.Any("values", packed_values))
				}

				if value_type == 3 {
					ptr := value_data.(int)
					value_data = string_table[ptr]
				}

				// entry_values.append(value_data)
				row[EntryName(entry_name)] = append(row[EntryName(entry_name)], value_data)
			}
			section = append(section, row)
		}
		result = append(result, Section{
			section_name: section_name,
			rows:         section,
		})

	}

	return result
}

type FoldValues bool

func Dump(path utils_types.FilePath) []string {
	bini := parse_file(path, FoldValues(false))

	var lines []string = make([]string, 0, 100)

	for _, section := range bini {

		// convert the entries in this section to strings and add to output
		lines = append(lines, fmt.Sprintf("[%s]", section.section_name))
		for _, row := range section.rows {
			// form key value pairs for each entry value. Expand tuples to remove quotes and brackets

			for row_key, row_values := range row {
				var formatted_values []string
				for _, value := range row_values {
					formatted_values = append(formatted_values, fmt.Sprintf("%v", value))
				}
				lines = append(lines,
					fmt.Sprintf("%s = %s", row_key, strings.Join(formatted_values, ", ")),
				)
			}

		}
		lines = append(lines, "") // add a blank line after each section

	}

	return lines
}

func IsBini(filepath utils_types.FilePath) bool {
	f, err := os.Open(filepath.ToString())
	logus.Log.CheckPanic(err,
		"file not founs in isBini check", utils_logus.FilePath(filepath))
	defer f.Close()

	bytes := make([]byte, 4)

	bufr := bufio.NewReader(f)
	bufr.Read(bytes)

	return string(bytes) == "BINI"
}
