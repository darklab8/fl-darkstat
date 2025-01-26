package exe_mapped

/*
This file is an adaptation of LibreLancer's DLL processing code.
The original C# code is licensed under the MIT License and available here:
https://github.com/Librelancer/Librelancer/blob/main/src/LibreLancer.Data/Dll/ResourceDll.cs
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"os"
	"strings"
	"unicode/utf16"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type IMAGE_RESOURCE_DIRECTORY struct {
	Characteristics      uint32
	TimeDateStamp        uint32
	MajorVersion         uint16
	MinorVersion         uint16
	NumberOfNamedEntries uint16
	NumberOfIdEntries    uint16
}

type IMAGE_RESOURCE_DIRECTORY_ENTRY struct {
	Name         uint32
	OffsetToData uint32
}

type IMAGE_RESOURCE_DATA_ENTRY struct {
	OffsetToData uint32
	Size         uint32
	CodePage     uint32
	Reserved     uint32
}

type ResourceTable struct {
	Type      uint32
	Resources []Resource
}

type Resource struct {
	Name    uint32
	Locales []ResourceData
}

type ResourceData struct {
	Locale uint32
	Data   []byte
}

const (
	RT_RCDATA                        = 23
	RT_STRING                        = 6
	IMAGE_RESOURCE_NAME_IS_STRING    = 0x80000000
	IMAGE_RESOURCE_DATA_IS_DIRECTORY = 0x80000000
)

type ResourceDll struct {
	Strings     map[int]string
	Infocards   map[int]string
	Dialogs     []BinaryResource
	Menus       []BinaryResource
	VersionInfo *VersionInfoResource
	SavePath    string
}

type BinaryResource struct {
	Name uint32
	Data []byte
}

type VersionInfoResource struct {
	Data []byte
}

func ParseDLL(fileData []byte, out *infocard.Config, globalOffset int) {
	rsrcOffset, rsrc, err := ReadPE(fileData)
	if err != nil {
		panic(err)
	}

	directory, err := Struct[IMAGE_RESOURCE_DIRECTORY](rsrc, 0)
	if err != nil {
		panic(err)
	}
	var resources []ResourceTable
	for i := 0; i < int(directory.NumberOfNamedEntries+directory.NumberOfIdEntries); i++ {
		off := 16 + (i * 8)
		entry, err := Struct[IMAGE_RESOURCE_DIRECTORY_ENTRY](rsrc, off)
		if err != nil {
			panic(err)
		}
		if (IMAGE_RESOURCE_NAME_IS_STRING & entry.Name) == IMAGE_RESOURCE_NAME_IS_STRING {
			continue
		}
		table, err := ReadResourceTable(rsrcOffset, DirOffset(entry.OffsetToData), rsrc, entry.Name)
		if err != nil {
			panic(err)
		}
		resources = append(resources, table)
	}

	for _, table := range resources {
		switch table.Type {
		case RT_RCDATA:
			for _, res := range table.Resources {
				idsId := globalOffset + int(res.Name)
				idx := 0
				count := len(res.Locales[0].Data)
				if count > 2 {
					if count%2 == 1 && res.Locales[0].Data[count-1] == 0 {
						count--
					}
					if res.Locales[0].Data[0] == 0xFF && res.Locales[0].Data[1] == 0xFE {
						idx += 2
					}
				}
				str, err := decodeUnicode(res.Locales[0].Data[idx:count])
				if err != nil {
					logus.Log.Warn("Infocard corrupt, skipping.", typelog.Any("id", idsId), typelog.Any("error", err))
					continue
				}
				out.Infocards[idsId] = infocard.NewInfocard(str)
			}
		case RT_STRING:
			for _, res := range table.Resources {
				blockId := globalOffset + int(res.Name-1)*16
				seg := res.Locales[0].Data
				reader := bytes.NewReader(seg)
				for j := 0; j < 16; j++ {
					var length uint16
					binary.Read(reader, binary.LittleEndian, &length)
					length *= 2
					if length != 0 {
						bytes := make([]byte, length)
						reader.Read(bytes)
						idsId := blockId + j
						str, err := decodeUnicode(bytes)
						if err != nil {
							logus.Log.Warn("Infostring corrupt, skipping.", typelog.Any("id", idsId), typelog.Any("error", err))
							continue
						}
						out.Infonames[idsId] = infocard.Infoname(str)
					}
				}
			}
		}
	}
}

func DirOffset(a uint32) int {
	return int(a & 0x7FFFFFFF)
}

func ReadResourceTable(rsrcOffset uint32, offset int, rsrc []byte, rtype uint32) (ResourceTable, error) {
	directory, err := Struct[IMAGE_RESOURCE_DIRECTORY](rsrc, offset)
	if err != nil {
		return ResourceTable{}, err
	}
	table := ResourceTable{Type: rtype}
	for i := 0; i < int(directory.NumberOfNamedEntries+directory.NumberOfIdEntries); i++ {
		off := offset + 16 + (i * 8)
		entry, err := Struct[IMAGE_RESOURCE_DIRECTORY_ENTRY](rsrc, off)
		if err != nil {
			return ResourceTable{}, err
		}
		res := Resource{Name: entry.Name}
		if (IMAGE_RESOURCE_DATA_IS_DIRECTORY & entry.OffsetToData) == IMAGE_RESOURCE_DATA_IS_DIRECTORY {
			langDirectory, err := Struct[IMAGE_RESOURCE_DIRECTORY](rsrc, DirOffset(entry.OffsetToData))
			if err != nil {
				return ResourceTable{}, err
			}
			for j := 0; j < int(langDirectory.NumberOfIdEntries+langDirectory.NumberOfNamedEntries); j++ {
				langOff := DirOffset(entry.OffsetToData) + 16 + (j * 8)
				langEntry, err := Struct[IMAGE_RESOURCE_DIRECTORY_ENTRY](rsrc, langOff)
				if err != nil {
					return ResourceTable{}, err
				}
				if (IMAGE_RESOURCE_DATA_IS_DIRECTORY & langEntry.OffsetToData) == IMAGE_RESOURCE_DATA_IS_DIRECTORY {
					return ResourceTable{}, errors.New("malformed .rsrc section")
				}
				dataEntry, err := Struct[IMAGE_RESOURCE_DATA_ENTRY](rsrc, int(langEntry.OffsetToData))
				if err != nil {
					return ResourceTable{}, err
				}
				dat := rsrc[(dataEntry.OffsetToData - rsrcOffset):(dataEntry.OffsetToData - rsrcOffset + dataEntry.Size)]
				res.Locales = append(res.Locales, ResourceData{Locale: langEntry.Name, Data: dat})
			}
		} else {
			return ResourceTable{}, errors.New("malformed .rsrc section")
		}
		table.Resources = append(table.Resources, res)
	}
	return table, nil
}

func Struct[T any](rawData []byte, offset int) (T, error) {
	var data T
	buf := bytes.NewReader(rawData[offset:])
	err := binary.Read(buf, binary.LittleEndian, &data)
	return data, err
}

func ReadPE(fullImage []byte) (uint32, []byte, error) {
	// Read the DOS header
	peOffset := binary.LittleEndian.Uint32(fullImage[60:64])
	peHeader := fullImage[peOffset:]

	// Read the PE header
	numberOfSections := binary.LittleEndian.Uint16(peHeader[6:8])
	sizeOfOptionalHeader := binary.LittleEndian.Uint16(peHeader[20:22])
	sectionTable := peHeader[24+sizeOfOptionalHeader:]

	// Find the resource section
	var rawStart, offset uint32
	for i := 0; i < int(numberOfSections); i++ {
		section := sectionTable[i*40 : (i+1)*40]
		sectionName := string(bytes.Trim(section[:8], "\x00"))
		if sectionName == ".rsrc" {
			offset = binary.LittleEndian.Uint32(section[12:16])
			// rsrcSize = binary.LittleEndian.Uint32(section[16:20])
			rawStart = binary.LittleEndian.Uint32(section[20:24])
			break
		}
	}

	if rawStart == 0 {
		logus.Log.Error("Resource section not found")
		return 0, nil, errors.New("resource section not found")
	}

	array := fullImage[rawStart:]
	return offset, array, nil
}

func decodeUnicode(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", errors.New("invalid byte array length for Unicode string")
	}
	u16 := make([]uint16, len(b)/2)
	for i := 0; i < len(u16); i++ {
		u16[i] = binary.LittleEndian.Uint16(b[i*2:])
	}
	runes := utf16.Decode(u16)
	return string(runes), nil
}

func ParseDLLs(dll_fnames []*file.File) *infocard.Config {
	out := infocard.NewConfig()

	for idx, name := range dll_fnames {
		data, err := os.ReadFile(name.GetFilepath().ToString())

		if logus.Log.CheckError(err, "unable to read dll") {
			continue
		}

		// if you inject "resources.dll" as 0 element of the list to process
		// despite it being not present in freelancer.ini and original Alex parsing script
		// then we go with global_offset from (idx) instead of (idx+1) as Alex had.
		global_offset := int(math.Pow(2, 16)) * (idx)

		func() {
			defer func() {
				if r := recover(); r != nil {
					logus.Log.Error("unable to read dll. Recovered by skipping dll.", typelog.String("filepath", name.GetFilepath().ToString()), typelog.Any("recover", r))
					if configs_settings.Env.Strict {
						panic(r)
					}
				}
			}()
			ParseDLL(data, out, global_offset)
		}()

	}

	return out
}

func GetAllInfocards(filesystem *filefind.Filesystem, dll_names []string) *infocard.Config {

	var files []*file.File
	for _, filename := range dll_names {
		dll_file := filesystem.GetFile(utils_types.FilePath(strings.ToLower(filename)))
		files = append(files, dll_file)
	}

	return ParseDLLs(files)
}
