package internal

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type FlagPointer struct {
	Required bool
	N        string
	Name     string
	Usage    string
	Type     reflect.Type
	Pointer  unsafe.Pointer
}

func (ptr *FlagPointer) FullName() string {
	sb := new(strings.Builder)

	if ptr.N != "" {
		sb.WriteString(ptr.N)
	} else if ptr.Name != "" {
		sb.WriteString(ptr.Name)
		return sb.String()
	}

	if ptr.Name != "" {
		sb.WriteString(", ")
		sb.WriteString(ptr.Name)
	}

	return sb.String()
}

func (ptr *FlagPointer) Parse(arg string) {
	switch {
	case parse[bool](arg, ptr, strconv.ParseBool):
		break
	case parse[string](arg, ptr, func(s string) (string, error) {
		return arg, nil
	}):
		break
	case parse[int](arg, ptr, func(s string) (int, error) {
		i, err := strconv.ParseInt(s, 0, 32)
		return int(i), err
	}):
		break
	case parse[int16](arg, ptr, func(s string) (int16, error) {
		i, err := strconv.ParseInt(s, 0, 16)
		return int16(i), err
	}):
		break
	case parse[int32](arg, ptr, func(s string) (int32, error) {
		i, err := strconv.ParseInt(s, 0, 32)
		return int32(i), err
	}):
		break
	case parse[int64](arg, ptr, func(s string) (int64, error) {
		return strconv.ParseInt(s, 0, 64)
	}):
		break
	case parse[uint](arg, ptr, func(s string) (uint, error) {
		u, err := strconv.ParseUint(s, 0, 32)
		return uint(u), err
	}):
		break
	case parse[uint16](arg, ptr, func(s string) (uint16, error) {
		u, err := strconv.ParseUint(s, 0, 16)
		return uint16(u), err
	}):
		break
	case parse[uint32](arg, ptr, func(s string) (uint32, error) {
		u, err := strconv.ParseUint(s, 0, 32)
		return uint32(u), err
	}):
		break
	case parse[uint64](arg, ptr, func(s string) (uint64, error) {
		u, err := strconv.ParseUint(s, 0, 64)
		return u, err
	}):
		break
	case parse[float32](arg, ptr, func(s string) (float32, error) {
		u, err := strconv.ParseFloat(s, 32)
		return float32(u), err
	}):
		break
	case parse[float64](arg, ptr, func(s string) (float64, error) {
		u, err := strconv.ParseFloat(s, 64)
		return u, err
	}):
		break
	case parse[float64](arg, ptr, func(s string) (float64, error) {
		u, err := strconv.ParseFloat(s, 64)
		return u, err
	}):
		break
	case parse[time.Duration](arg, ptr, time.ParseDuration):
		break
	default:
		fmt.Println("Invalid argument for " + ptr.Type.Name() + " flag " + ptr.FullName())
		os.Exit(0)
	}
}

func (ptr *FlagPointer) IsBoolean() bool {
	return reflect.TypeFor[bool]().AssignableTo(ptr.Type)
}

func parse[T any](arg string, ptr *FlagPointer, parser func(string) (T, error)) bool {
	if t := reflect.TypeFor[T](); !ptr.Type.AssignableTo(t) {
		return false
	}

	if parsed, err := parser(arg); err == nil {
		*((*T)(ptr.Pointer)) = parsed
		return true
	}

	return false
}
