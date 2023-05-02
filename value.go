package argsx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Value struct {
	fullkey string
	payload string
}

// parser is a generic type convert string to T
type parser[T any] func(payload string) (T, error)

// NewValue returns Value instance with payload
func NewValue(payload string) Value {
	return Value{payload: payload}
}

// get returns parse result of T type, if payload is not specified return default value or zero value
func get[T any](v Value, dv []T, parse parser[T]) (t T, err error) {
	if len(v.payload) == 0 {
		if len(dv) > 0 {
			return dv[0], nil
		}
		if len(v.fullkey) == 0 {
			return t, errors.New("invalid value: empty")
		}
		return t, fmt.Errorf("args not specified value for key: `%s`", v.fullkey)
	}
	return parse(v.payload)
}

// must check the err if nil then return val otherwise return zero value of T type
func must[T any](val T, err error) (t T) {
	if err != nil {
		return
	}
	return val
}

// toSlice convert payload string to slice of T type
func toSlice[T any](payload, delimiter string, parse parser[T], dv ...T) ([]T, error) {
	var slice []T
	arr := strings.Split(payload, delimiter)
	for _, str := range arr {
		if len(str) == 0 {
			if len(dv) > 0 {
				slice = append(slice, dv[0])
			}
			continue
		}

		b, err := parse(str)
		if err != nil {
			return nil, err
		}
		slice = append(slice, b)
	}
	return slice, nil
}

// String returns string value
//
//	NewValue("string value").String() // "string value", nil
//	NewValue("").String() // "", error
//	NewValue("").String("default value") // "default value", nil
func (v Value) String(dv ...string) (string, error) {
	return get(v, dv, func(payload string) (string, error) {
		return payload, nil
	})
}

// MustString returns string value ignore error
//
//	NewValue("must string").MustString() // "must string"
//	NewValue("").MustString() // ""
//	NewValue("").MustString("default must string") // "default must string"
func (v Value) MustString(dv ...string) string {
	return must(v.String(dv...))
}

// StringSlice returns []string and error
//
//	NewValue("A,B,C").StringSlice() // []string{"A", "B", "C"}, nil
//	NewValue("").StringSlice() // nil, error
//	NewValue("").StringSlice(WithDefault[string]("D", "E", "F")) // []string{"D", "E", "F"}, error
//	NewValue("G/H/I").StringSlice(WithDelimiter[string]("/")) // []string{"G", "H", "I"}, error
func (v Value) StringSlice(opts ...Option[string]) ([]string, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]string, error) {
		return strings.Split(payload, option.delimiter), nil
	})
}

// MustStringSlice returns []string ignore error
//
//	NewValue("A,B,C").MustStringSlice() // []string{"A", "B", "C"}
//	NewValue("").MustStringSlice() // nil
//	NewValue("").MustStringSlice(WithDefault[string]("D", "E", "F")) // []string{"D", "E", "F"}
//	NewValue("G.H.I").MustStringSlice(WithDelimiter[string](".")) // []string{"G", "H", "I"}
func (v Value) MustStringSlice(opts ...Option[string]) []string {
	return must(v.StringSlice(opts...))
}

// Bool returns bool value of payload
//
//	NewValue("true").Bool() // true, nil
//	NewValue("false").Bool() // false, nil
//	NewValue("").Bool() // true, nil
//	NewValue("").Bool(false) // false, nil
//	NewValue("abc").Bool() // false, error
func (v Value) Bool(dv ...bool) (bool, error) {
	return get[bool](v, append(dv, true), strconv.ParseBool)
}

// MustBool returns bool value of payload ignore error
//
//	NewValue("true").MustBool() // true
//	NewValue("false").MustBool() // false
//	NewValue("").MustBool() // true
//	NewValue("").MustBool(false) // false
//	NewValue("abc").MustBool() // false
func (v Value) MustBool(dv ...bool) bool {
	return must(v.Bool(dv...))
}

// BoolSlice returns []bool
//
//	NewValue("true,T,True,1").BoolSlice() // []bool{true, true, true, true}, nil
//	NewValue("").BoolSlice() // nil, error
//	NewValue("true,F").BoolSlice(WithDelimiter[bool](";")) // []bool{true, false}, nil
//	NewValue("").BoolSlice(WithDefault[bool](true, false)) // []bool{true, false}, nil
func (v Value) BoolSlice(opts ...Option[bool]) ([]bool, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]bool, error) {
		return toSlice(payload, option.delimiter, strconv.ParseBool, true)
	})
}

// MustBoolSlice return []bool if error not nil will be ignored
//
//	NewValue("true,T,True,1").MustBoolSlice() // []bool{true, true, true, true}
//	NewValue("").MustBoolSlice() // nil
//	NewValue("true,F").MustBoolSlice(WithDelimiter[bool](";")) // []bool{true, false}
//	NewValue("").MustBoolSlice(WithDefault[bool](true, false)) // []bool{true, false}
func (v Value) MustBoolSlice(opts ...Option[bool]) []bool {
	return must(v.BoolSlice(opts...))
}

// Int returns int value
//
//	NewValue("5").Int() // 5, nil
//	NewValue("").Int() // 0, error
//	NewValue("").Int(7) // 7, nil
//	NewValue("a").Int() // 0, error
func (v Value) Int(dv ...int) (int, error) {
	return get(v, dv, strconv.Atoi)
}

// MustInt returns int value or zero when payload not a valid number ignore error
//
//	NewValue("5").MustInt() // 5
//	NewValue("").MustInt() // 0
//	NewValue("").MustInt(7) // 7
//	NewValue("a").MustInt() // 0
func (v Value) MustInt(dv ...int) int {
	return must(v.Int(dv...))
}

// IntSlice return []int and error
//
//	NewValue("1,2,3").IntSlice() // []int{1, 2, 3}, nil
//	NewValue("").IntSlice() // nil, error
//	NewValue("0,7,a,b").IntSlice() // nil, error
//	NewValue("").IntSlice(WithDefault[int](4, 5, 6)) // []int{4, 5, 6}, nil
//	NewValue("7;8;9").IntSlice(WithDelimiter[int](";")) // []int{7, 8, 9}, nil
func (v Value) IntSlice(opts ...Option[int]) ([]int, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]int, error) {
		return toSlice(payload, option.delimiter, strconv.Atoi)
	})
}

// MustIntSlice returns []int if error is not nil will be ignored
//
//	NewValue("1,2,3").MustIntSlice() // []int{1, 2, 3}, nil
//	NewValue("").MustIntSlice() // nil
//	NewValue("0,7,a,b").MustIntSlice() // nil
//	NewValue("").MustIntSlice(WithDefault[int](4, 5, 6)) // []int{4, 5, 6}
//	NewValue("7;8;9").MustIntSlice(WithDelimiter[int](7, 8, 9)) // []int{7, 8, 9}
func (v Value) MustIntSlice(opts ...Option[int]) []int {
	return must(v.IntSlice(opts...))
}

// parseInt8 parse payload string to int8
func parseInt8(payload string) (int8, error) {
	if i, err := strconv.ParseInt(payload, 0, 8); err == nil {
		return int8(i), nil
	} else {
		return 0, err
	}
}

// Int8 returns int8
//
//	NewValue("9").Int8() // int8(9), nil
//	NewValue("").Int8() // int8(0), error
//	NewValue("").Int8(7) // int8(7), nil
//	NewValue("a").Int8() // int8(0), error
func (v Value) Int8(dv ...int8) (int8, error) {
	return get(v, dv, parseInt8)
}

// MustInt8 returns int8 if error not nil will be ignored
//
//	NewValue("9").MustInt8() // int8(9)
//	NewValue("").MustInt8() // int8(0)
//	NewValue("").MustInt8(7) // int8(7)
//	NewValue("a").MustInt8() // int8(0)
func (v Value) MustInt8(dv ...int8) int8 {
	return must(v.Int8(dv...))
}

// Int8Slice returns []int8
//
//	NewValue("7,8,9").Int8Slice() // []int8{7, 8, 9}, nil
//	NewValue("").Int8Slice() // nil, error
//	NewValue("").Int8Slice(WithDefault[int8](4, 5, 6)) // []int8{4, 5, 6}, nil
//	NewValue("1;2;3").Int8Slice(WithDelimiter[int8](";")) // []int8{1, 2, 3}, nil
func (v Value) Int8Slice(opts ...Option[int8]) ([]int8, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]int8, error) {
		return toSlice(payload, option.delimiter, parseInt8)
	})
}

// MustInt8Slice returns []int8 if error not nil will be ignored
//
//	NewValue("7,8,9").MustInt8Slice() // []int8{7, 8, 9}
//	NewValue("").MustInt8Slice() // nil
//	NewValue("").MustInt8Slice(WithDefault[int8](4, 5, 6)) // []int8{4, 5, 6}
//	NewValue("1;2;3").MustInt8Slice(WithDelimiter[int8](";")) // []int8{1, 2, 3}
func (v Value) MustInt8Slice(opts ...Option[int8]) []int8 {
	return must(v.Int8Slice(opts...))
}

// parseInt16 returns int16
func parseInt16(payload string) (int16, error) {
	if i, err := strconv.ParseInt(payload, 0, 16); err == nil {
		return int16(i), nil
	} else {
		return 0, err
	}
}

// Int16 returns int16
//
//	NewValue("16").Int16() // int16(16), nil
//	NewValue("").Int16() // int16(0), error
//	NewValue("").Int16(7) // int16(7), nil
//	NewValue("a").Int16() // int16(0), error
func (v Value) Int16(dv ...int16) (int16, error) {
	return get(v, dv, parseInt16)
}

// MustInt16 returns int16 if error is not nil will be ignored
//
//	NewValue("16").MustInt16() // int16(16)
//	NewValue("").MustInt16() // int16(0)
//	NewValue("").MustInt16(7) // int16(7)
//	NewValue("a").MustInt16() // int16(0)
func (v Value) MustInt16(dv ...int16) int16 {
	return must(v.Int16(dv...))
}

// Int16Slice returns []int16
//
//	NewValue("7,8,9").Int16Slice() // []int16{7, 8, 9}, nil
//	NewValue("").Int16Slice() // nil, error
//	NewValue("").Int16Slice(WithDefault[int16](4, 5, 6)) // []int16{4, 5, 6}, nil
//	NewValue("1;2;3").Int16Slice(WithDelimiter[int16](";")) // []int16{1, 2, 3}, nil
func (v Value) Int16Slice(opts ...Option[int16]) ([]int16, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]int16, error) {
		return toSlice(payload, option.delimiter, parseInt16)
	})
}

// MustInt16Slice returns []int16 if error not nil will be ignored
//
//	NewValue("7,8,9").MustInt16Slice() // []int16{7, 8, 9}
//	NewValue("").MustInt16Slice() // nil
//	NewValue("").MustInt16Slice(WithDefault[int16](4, 5, 6)) // []int16{4, 5, 6}
//	NewValue("1;2;3").MustInt16Slice(WithDelimiter[int16](";")) // []int16{1, 2, 3}
func (v Value) MustInt16Slice(opts ...Option[int16]) []int16 {
	return must(v.Int16Slice(opts...))
}

// parseInt32 returns int32 value
func parseInt32(payload string) (int32, error) {
	if i, err := strconv.ParseInt(payload, 0, 32); err == nil {
		return int32(i), nil
	} else {
		return 0, err
	}
}

// Int32 returns int32
//
//	NewValue("16").Int32() // int32(16), nil
//	NewValue("").Int32() // int32(0), error
//	NewValue("").Int32(7) // int32(7), nil
//	NewValue("a").Int32() // int32(0), error
func (v Value) Int32(dv ...int32) (int32, error) {
	return get(v, dv, parseInt32)
}

// MustInt32 returns int32 if error is not nil will be ignored
//
//	NewValue("16").MustInt32() // int32(16)
//	NewValue("").MustInt32() // int32(0)
//	NewValue("").MustInt32(7) // int32(7)
//	NewValue("a").MustInt32() // int32(0)
func (v Value) MustInt32(dv ...int32) int32 {
	return must(v.Int32(dv...))
}

// Int32Slice returns []int32
//
//	NewValue("7,8,9").Int32Slice() // []int32{7, 8, 9}, nil
//	NewValue("").Int32Slice() // nil, error
//	NewValue("").Int32Slice(WithDefault[int32](4, 5, 6)) // []int32{4, 5, 6}, nil
//	NewValue("1;2;3").Int32Slice(WithDelimiter[int32](";")) // []int32{1, 2, 3}, nil
func (v Value) Int32Slice(opts ...Option[int32]) ([]int32, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]int32, error) {
		return toSlice(payload, option.delimiter, parseInt32)
	})
}

// MustInt32Slice returns []int32 if error not nil will be ignored
//
//	NewValue("7,8,9").MustInt32Slice() // []int32{7, 8, 9}
//	NewValue("").MustInt32Slice() // nil
//	NewValue("").MustInt32Slice(WithDefault[int32](4, 5, 6)) // []int32{4, 5, 6}
//	NewValue("1;2;3").MustInt32Slice(WithDelimiter[int32](";")) // []int32{1, 2, 3}
func (v Value) MustInt32Slice(opts ...Option[int32]) []int32 {
	return must(v.Int32Slice(opts...))
}

// parseInt64 returns int64 value
func parseInt64(payload string) (int64, error) {
	return strconv.ParseInt(payload, 0, 64)
}

// Int64 returns int64
//
//	NewValue("16").Int64() // int64(16), nil
//	NewValue("").Int64() // int64(0), error
//	NewValue("").Int64(7) // int64(7), nil
//	NewValue("a").Int64() // int64(0), error
func (v Value) Int64(dv ...int64) (int64, error) {
	return get(v, dv, parseInt64)
}

// MustInt64 returns int64 if error is not nil will be ignored
//
//	NewValue("16").MustInt64() // int64(16)
//	NewValue("").MustInt64() // int64(0)
//	NewValue("").MustInt64(7) // int64(7)
//	NewValue("a").MustInt64() // int64(0)
func (v Value) MustInt64(dv ...int64) int64 {
	return must(v.Int64(dv...))
}

// Int64Slice returns []int64
//
//	NewValue("7,8,9").Int64Slice() // []int64{7, 8, 9}, nil
//	NewValue("").Int64Slice() // nil, error
//	NewValue("").Int64Slice(WithDefault[int64](4, 5, 6)) // []int64{4, 5, 6}, nil
//	NewValue("1;2;3").Int64Slice(WithDelimiter[int64](";")) // []int64{1, 2, 3}, nil
func (v Value) Int64Slice(opts ...Option[int64]) ([]int64, error) {
	option := getOpts(opts)
	return get(v, option.getDefault(), func(payload string) ([]int64, error) {
		return toSlice(payload, option.delimiter, parseInt64)
	})
}

// MustInt64Slice returns []int64 if error not nil will be ignored
//
//	NewValue("7,8,9").MustInt64Slice() // []int64{7, 8, 9}
//	NewValue("").MustInt64Slice() // nil
//	NewValue("").MustInt64Slice(WithDefault[int64](4, 5, 6)) // []int64{4, 5, 6}
//	NewValue("1;2;3").MustInt64Slice(WithDelimiter[int64](";")) // []int64{1, 2, 3}
func (v Value) MustInt64Slice(opts ...Option[int64]) []int64 {
	return must(v.Int64Slice(opts...))
}
