package flagg

import (
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/raiiga/flagg/internal"
)

const (
	tag    = "flagg"
	colon  = ":"
	empty  = ""
	cutset = ", "
	escape = "\\"
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

	if fileInfo, _ := os.Stdin.Stat(); fileInfo.Mode()&os.ModeNamedPipe != 0 {
		return m.Parser.ParseWithPipe(args, os.Stdin)
	}

	return m.Parser.Parse(args)
}

func (m *flagg) process(lookup string, fieldValue reflect.Value) error {
	params := map[string]string{}
	split := regexp.MustCompile(`.*?[^\\](,|$)`).FindAllString(lookup, -1)

	for _, s := range split {
		kv := strings.Trim(s, cutset)

		if i := strings.SplitN(kv, colon, 2); len(i) == 2 {
			params[strings.TrimSpace(i[0])] = strings.ReplaceAll(strings.TrimSpace(i[1]), escape, empty)
		}
	}

	return m.Parser.FromMap(fieldValue, params)
}
