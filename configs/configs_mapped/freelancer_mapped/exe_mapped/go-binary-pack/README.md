# Go BinaryPack


BinaryPack is a simple Golang library which implements some functionality of Python's [struct](https://docs.python.org/2/library/struct.html) package.

This one is forked from [BinaryPack](https://github.com/roman-kachanovsky/go-binary-pack) with more supported `pad byte` type

**Format characters**

Format characters (some characters like H have been reserved for future implementation of unsigned numbers):

```
    ? - bool, packed size 1 byte
    b - int8, packed size 1 bytes
    B - uint8, packed size 1 bytes
    h - int16, packed size 2 bytes
    H - uint16, packed size 2 bytes
    i, l - int32, packed size 4 bytes
    I, L - int32, packed size 4 bytes
    q - int64, packed size 8 bytes
    Q - uint64, packed size 8 bytes
    f - float32, packed size 4 bytes
    d - float64, packed size 8 bytes
    Ns - string, packed size N bytes, N is a number of runes to pack/unpack
```


**Install**

`go get github.com/roman-kachanovsky/go-binary-pack/binary-pack`

**How to use**

```go
// Prepare format (slice of strings)
format := []string{"I", "?", "d", "6s"}

// Prepare values to pack
values := []interface{}{4, true, 3.14, "Golang"}

// Create BinaryPack object
bp := new(BinaryPack)

// Pack values to []byte
data, err := bp.Pack(format, values)

// Unpack binary data to []interface{}
unpacked_values, err := bp.UnPack(format, data)

// You can calculate size of expected binary data by format
size, err := bp.CalcSize(format)

```
