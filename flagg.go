package flagg

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type flagg struct {
	Parser *parser
}

func New(name string) *flagg {
	var sb = new(strings.Builder)
	sb.WriteString(fmt.Sprintf("Usage of %s:", name))

	return &flagg{
		Parser: &parser{
			Usage:   sb,
			FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
		},
	}
}

func (m *flagg) Map(entity any) (bool, error) {
	return m.MapFromArgs(entity, os.Args[1:])
}

func (m *flagg) MapFromArgs(entity any, args []string) (bool, error) {
	var (
		typeOf  = reflect.TypeOf(entity).Elem()
		valueOf = reflect.ValueOf(entity).Elem()
	)

	for i, l := 0, typeOf.NumField(); i < l; i++ {
		if lookup, ok := typeOf.Field(i).Tag.Lookup("flagg"); ok {
			m.process(lookup, valueOf.Field(i))
		}
	}

	return m.Parser.Parse(args)
}

func (m *flagg) process(lookup string, fieldValue reflect.Value) {
	var (
		params = map[string]string{}
		split  = strings.Split(lookup, ",")
	)

	for _, s := range split {
		if i := strings.Split(strings.TrimSpace(s), ":"); len(i) == 2 {
			params[strings.TrimSpace(i[0])] = strings.TrimSpace(i[1])
		}
	}

	m.Parser.FromMap(fieldValue, params)
}
