package flagg

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type mapper struct {
	FlagSet *flag.FlagSet
}

func NewMapper(name string, handling flag.ErrorHandling) *mapper {
	return &mapper{FlagSet: flag.NewFlagSet(name, handling)}
}

func (m *mapper) Map(entity any) error {
	return m.MapFromArgs(entity, os.Args[1:])
}

func (m *mapper) MapFromArgs(entity any, args []string) error {
	typeOf, valueOf := reflect.TypeOf(entity).Elem(), reflect.ValueOf(entity).Elem()

	for i, l := 0, typeOf.NumField(); i < l; i++ {
		fieldTyp, fieldVal := typeOf.Field(i), valueOf.Field(i)

		usage, _ := fieldTyp.Tag.Lookup(usageTagName)
		value, _ := fieldTyp.Tag.Lookup(valueTagName)

		if full, f := fieldTyp.Tag.Lookup(longTagName); f {
			m.processFlag(full, usage, value, fieldVal)
		}

		if short, s := fieldTyp.Tag.Lookup(shortTagName); s {
			m.processFlag(short, usage, value, fieldVal)
		}
	}

	return m.FlagSet.Parse(args)
}

func (m *mapper) processFlag(name, usage, value string, fieldVal reflect.Value) {
	ptr, typ := unsafe.Pointer(fieldVal.UnsafeAddr()), fieldVal.Type()

	if t := reflect.TypeFor[bool](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseBool(value)
		m.FlagSet.BoolVar((*bool)(ptr), name, parsed, usage)
		return
	}

	if t := reflect.TypeFor[string](); typ.AssignableTo(t) {
		m.FlagSet.StringVar((*string)(ptr), name, value, usage)
		return
	}

	if t := reflect.TypeFor[int](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseInt(value, 0, 32)
		m.FlagSet.IntVar((*int)(ptr), name, int(parsed), usage)
		return
	}

	if t := reflect.TypeFor[int64](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseInt(value, 0, 64)
		m.FlagSet.Int64Var((*int64)(ptr), name, parsed, usage)
		return
	}

	if t := reflect.TypeFor[uint](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseUint(value, 0, 32)
		m.FlagSet.UintVar((*uint)(ptr), name, uint(parsed), usage)
		return
	}

	if t := reflect.TypeFor[uint64](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseUint(value, 0, 64)
		m.FlagSet.Uint64Var((*uint64)(ptr), name, parsed, usage)
		return
	}

	if t := reflect.TypeFor[float64](); typ.AssignableTo(t) {
		parsed, _ := strconv.ParseFloat(value, 64)
		m.FlagSet.Float64Var((*float64)(ptr), name, parsed, usage)
		return
	}

	if t := reflect.TypeFor[time.Duration](); typ.AssignableTo(t) {
		parsed, _ := time.ParseDuration(value)
		m.FlagSet.DurationVar((*time.Duration)(ptr), name, parsed, usage)
		return
	}

	if t := reflect.TypeFor[func(string) error](); typ.AssignableTo(t) {
		m.FlagSet.Func(name, usage, *(*func(string) error)(ptr))
		return
	}

	if t := reflect.TypeFor[func() error](); typ.AssignableTo(t) {
		m.FlagSet.BoolFunc(name, usage, func(string) error { return (*(*func() error)(ptr))() })
		return
	}
}
