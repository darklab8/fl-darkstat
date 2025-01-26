package bin

/*
Code simplifying interactions with go-binary-pack
*/

import (
	"bytes"

	gbp "github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped/go-binary-pack"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

type Bdatas struct {
	bdatas map[int][]byte
}

func NewBDatas() *Bdatas {
	b := &Bdatas{
		bdatas: make(map[int][]byte),
	}
	b.bdatas[1] = make([]byte, 1)
	b.bdatas[2] = make([]byte, 2)
	b.bdatas[3] = make([]byte, 3)
	b.bdatas[4] = make([]byte, 4)
	b.bdatas[8] = make([]byte, 8)
	return b
}

func (b *Bdatas) GetBData(BytesToRead int) []byte {
	var bdata []byte = b.bdatas[BytesToRead]
	for i := range bdata {
		bdata[i] = 0
	}
	return bdata
}

var packer = new(gbp.BinaryPack)

func Read(
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) ([]interface{}, int, error) {
	returned_n, err := fh.Read(byte_data)

	if err != nil {
		var UnpackErrValue []interface{}
		return UnpackErrValue, returned_n, err
	}

	unpacked_value, err := packer.UnPack(format, byte_data)

	if err != nil {
		var UnpackErrValue []interface{}
		return UnpackErrValue, returned_n, err
	}
	return unpacked_value, returned_n, nil
}

func Read2(
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) ([]interface{}, int) {
	// version withtout err check
	values, returned_n, err := Read(fh, byte_data, format)
	logus.Log.CheckPanic(err, "failed to read")
	return values, returned_n
}

func Unpack[returnType any](
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) (returnType, int, error) {
	packed_values, return_n, err := Read(fh, byte_data, format)
	value := packed_values[0].(returnType)
	return value, return_n, err
}

func Unpack2[returnType any](
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) (returnType, int) {
	value, return_n, err := Unpack[returnType](fh, byte_data, format)
	logus.Log.CheckPanic(err, "failed to read unpack")
	return value, return_n
}

func Unpack3[returnType any](
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) returnType {
	value, _, err := Unpack[returnType](fh, byte_data, format)
	logus.Log.CheckPanic(err, "failed to read unpack")
	return value
}
