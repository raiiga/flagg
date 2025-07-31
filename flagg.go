package flagg

import (
	"github.com/raiiga/flagg/internal"
	"os"
	"reflect"
	"strings"
)

const (
	tag   = "flagg"
	colon = ":"
	comma = ","
)

type flagg struct {
	Parser *internal.Parser
}

func New(usage string) *flagg {
	sb := new(strings.Builder)
	sb.WriteString(usage)

	return &flagg{
		Parser: &internal.Parser{
			Usage:       sb,
			Pointers:    make([]*internal.FlagPointer, 0),
			PointersMap: make(map[string]*internal.FlagPointer),
		},
	}
}

func (m *flagg) Map(entity any) error {
	return m.MapFromArgs(entity, os.Args[1:])
}

func (m *flagg) MapFromArgs(entity any, args []string) error {
	typeOf := reflect.TypeOf(entity).Elem()
	valueOf := reflect.ValueOf(entity).Elem()

	for i, l := 0, typeOf.NumField(); i < l; i++ {
		if lookup, ok := typeOf.Field(i).Tag.Lookup(tag); ok {
			if err := m.process(lookup, valueOf.Field(i)); err != nil {
				return err
			}
		}
	}

	if fileInfo, _ := os.Stdin.Stat(); fileInfo.Mode()&os.ModeCharDevice == 0 {
		return m.Parser.ParseWithPipe(args, os.Stdin)
	}

	return m.Parser.Parse(args)
}

func (m *flagg) process(lookup string, fieldValue reflect.Value) error {
	params := map[string]string{}
	split := strings.Split(lookup, comma)

	for _, s := range split {
		if i := strings.Split(strings.TrimSpace(s), colon); len(i) == 2 {
			params[strings.TrimSpace(i[0])] = strings.TrimSpace(i[1])
		}
	}

	return m.Parser.FromMap(fieldValue, params)
}
