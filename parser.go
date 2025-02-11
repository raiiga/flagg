package flagg

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	N       = "n"
	Name    = "name"
	value   = "value"
	Usage   = "usage"
	Default = "default"
)

var (
	lineSeparator = fmt.Sprintln()
)

type parser struct {
	FlagSet *flag.FlagSet
	Usage   *strings.Builder
}

func (p *parser) FromMap(fieldValue reflect.Value, parameterMap map[string]string) {
	if parameterMap[value] != "" {
		p.populateFlag(fieldValue, parameterMap[value], parameterMap[N], parameterMap[Name], parameterMap[Usage])
	} else {
		p.populateFlag(fieldValue, parameterMap[Default], parameterMap[N], parameterMap[Name], parameterMap[Usage])
	}
}

func (p *parser) populateFlag(fieldValue reflect.Value, flagValue, flagN, flagName, flagUsage string) {
	var (
		mask = byte(0)
		ptr  = unsafe.Pointer(fieldValue.UnsafeAddr())
		typ  = fieldValue.Type()
	)

	if flagN != "" {
		p.populateFlag_(ptr, typ, flagN, flagValue)
		mask |= 1
	}

	if flagName != "" {
		p.populateFlag_(ptr, typ, flagName, flagValue)
		mask |= 2
	}

	p.populateUsage(typ, flagN, flagName, flagUsage, mask)
}

func (p *parser) populateFlag_(pointer unsafe.Pointer, flagType reflect.Type, flagName, flagValue string) {
	if t := reflect.TypeFor[bool](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseBool(flagValue)
		p.FlagSet.BoolVar((*bool)(pointer), flagName, parsed, "")
		return
	}

	if t := reflect.TypeFor[string](); flagType.AssignableTo(t) {
		p.FlagSet.StringVar((*string)(pointer), flagName, flagValue, "")
		return
	}

	if t := reflect.TypeFor[int](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseInt(flagValue, 0, 32)
		p.FlagSet.IntVar((*int)(pointer), flagName, int(parsed), "")
		return
	}

	if t := reflect.TypeFor[int64](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseInt(flagValue, 0, 64)
		p.FlagSet.Int64Var((*int64)(pointer), flagName, parsed, "")
		return
	}

	if t := reflect.TypeFor[uint](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseUint(flagValue, 0, 32)
		p.FlagSet.UintVar((*uint)(pointer), flagName, uint(parsed), "")
		return
	}

	if t := reflect.TypeFor[uint64](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseUint(flagValue, 0, 64)
		p.FlagSet.Uint64Var((*uint64)(pointer), flagName, parsed, "")
		return
	}

	if t := reflect.TypeFor[float64](); flagType.AssignableTo(t) {
		parsed, _ := strconv.ParseFloat(flagValue, 64)
		p.FlagSet.Float64Var((*float64)(pointer), flagName, parsed, "")
		return
	}

	if t := reflect.TypeFor[time.Duration](); flagType.AssignableTo(t) {
		parsed, _ := time.ParseDuration(flagValue)
		p.FlagSet.DurationVar((*time.Duration)(pointer), flagName, parsed, "")
		return
	}

	if t := reflect.TypeFor[func() error](); flagType.AssignableTo(t) {
		once := sync.OnceValue(*(*func() error)(pointer))
		p.FlagSet.BoolFunc(flagName, "", func(string) error { return once() })
		return
	}

	if t := reflect.TypeFor[func(string) error](); flagType.AssignableTo(t) {
		mutex := new(sync.Mutex)
		p.FlagSet.Func(flagName, "", func(s string) error {
			if !mutex.TryLock() {
				return nil
			}
			return (*(*func(string) error)(pointer))(s)
		})
	}
}

func (p *parser) populateUsage(typ reflect.Type, flagN, flagName, flagUsage string, mask byte) {
	type_ := typ.String()
	if type_ == "bool" || strings.Contains(type_, "func() error") {
		type_ = ""
	} else if strings.Contains(type_, "func(string) error") {
		type_ = "string"
	}

	switch mask {
	case 1:
		p.Usage.WriteString(fmt.Sprintf("%s-%s %s", lineSeparator, flagN, type_))
	case 2:
		p.Usage.WriteString(fmt.Sprintf("%s-%s %s", lineSeparator, flagName, type_))
	case 3:
		p.Usage.WriteString(fmt.Sprintf("%s-%s, -%s %s", lineSeparator, flagN, flagName, type_))
	default:
		break
	}

	p.Usage.WriteString(fmt.Sprintf("%s\t%s", lineSeparator, flagUsage))
}

func (p *parser) Parse(args []string) (bool, error) {
	p.FlagSet.Usage = func() {
		println(p.Usage.String())
	}

	if err := p.FlagSet.Parse(args); errors.Is(err, flag.ErrHelp) {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return false, nil
}
