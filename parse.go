package config

import (
	"fmt"
	"strconv"
)

func parseInt(i *int, s string) (err error) {
	*i, err = strconv.Atoi(s)
	return
}

func parseInt64(i *int64, s string) (err error) {
	*i, err = strconv.ParseInt(s, 10, 64)
	return
}

func parseInt32(i *int32, s string) (err error) {
	var i64 int64
	i64, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	if i64 < int32Min || i64 > int32Max {
		err = fmt.Errorf("integer %d exceeds int32 size (%d to %d)", i64, int32Min, int32Max)
		return
	}
	*i = int32(i64)
	return err
}

func parseInt16(i *int16, s string) (err error) {
	var i64 int64
	i64, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	if i64 < int16Min || i64 > int16Max {
		err = fmt.Errorf("integer %d exceeds int16 size (%d to %d)", i64, int16Min, int16Max)
		return
	}
	*i = int16(i64)
	return err
}

func parseInt8(i *int8, s string) (err error) {
	var i64 int64
	i64, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	if i64 < int8Min || i64 > int8Max {
		err = fmt.Errorf("integer %d exceeds int8 size (%d to %d)", i64, int8Min, int8Max)
		return
	}
	*i = int8(i64)
	return err
}

func parseUint(i *uint, s string) (err error) {
	var u64 uint64
	u64, err = strconv.ParseUint(s, 10, intSize)
	if u64 > uintMax {
		err = fmt.Errorf("unsigned integer %d exceeds max uint (%d)", u64, uintMax)
		return
	}
	*i = uint(u64)
	return
}

func parseUint64(i *uint64, s string) (err error) {
	*i, err = strconv.ParseUint(s, 10, 64)
	return err
}

func parseUint32(i *uint32, s string) (err error) {
	var u64 uint64
	u64, err = strconv.ParseUint(s, 10, 64)
	if u64 > uintMax {
		err = fmt.Errorf("unsigned integer %d exceeds uint size (%d to %d)", u64, uintMin, uintMax)
		return
	}
	*i = uint32(u64)
	return
}

func parseUint16(i *uint16, s string) (err error) {
	var u64 uint64
	u64, err = strconv.ParseUint(s, 10, 64)
	if u64 > uintMax {
		err = fmt.Errorf("unsigned integer %d exceeds uint size (%d to %d)", u64, uintMin, uintMax)
		return
	}
	*i = uint16(u64)
	return
}

func parseUint8(i *uint8, s string) (err error) {
	var u64 uint64
	u64, err = strconv.ParseUint(s, 10, 64)
	if u64 > uintMax {
		err = fmt.Errorf("unsigned integer %d exceeds uint size (%d to %d)", u64, uintMin, uintMax)
		return
	}
	*i = uint8(u64)
	return
}

// MARK: Size Constants

const (
	intSize = strconv.IntSize

	// signed integers
	intMin   int64 = -1 << (intSize - 1)
	intMax   int64 = 1<<(intSize-1) - 1
	int64Min int64 = -1 << 63
	int64Max int64 = 1<<63 - 1
	int32Min int64 = -1 << 31
	int32Max int64 = 1<<31 - 1
	int16Min int64 = -1 << 15
	int16Max int64 = 1<<15 - 1
	int8Min  int64 = -1 << 7
	int8Max  int64 = 1<<7 - 1

	// unsigned integers
	uintMin   uint64 = 0
	uintMax   uint64 = 1<<intSize - 1
	uint64Max uint64 = 1<<64 - 1
	uint32Max uint64 = 1<<32 - 1
	uint16Max uint64 = 1<<16 - 1
	uint8Max  uint64 = 1<<8 - 1

	// floats
	float32Min = -float32Max
	float32Max = 0x1p127 * (1 + (1 - 0x1p-23))
	float64Min = -float64Max
	float64Max = 0x1p1023 * (1 + (1 - 0x1p-52))
)
