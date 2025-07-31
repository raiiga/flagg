package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

const (
	N            = "n"
	Name         = "name"
	Usage        = "usage"
	Required     = "required"
	Space        = "\u0020"
	LeftPadding  = 2
	RightPadding = 4
)

type Parser struct {
	Usage       *strings.Builder
	Pointers    []*FlagPointer
	PointersMap map[string]*FlagPointer
}

func (p *Parser) FromMap(fieldValue reflect.Value, parameterMap map[string]string) error {
	flagPtr := &FlagPointer{
		Usage:   parameterMap[Usage],
		Type:    fieldValue.Type(),
		Pointer: unsafe.Pointer(fieldValue.UnsafeAddr()),
	}

	flagPtr.Required, _ = strconv.ParseBool(parameterMap[Required])

	if parameterMap[N] != "" {
		flagPtr.N = "-" + parameterMap[N]

		if p.PointersMap[flagPtr.N] != nil {
			return errors.New("flag " + flagPtr.N + " has been already defined")
		}

		p.PointersMap[flagPtr.N] = flagPtr
	}

	if parameterMap[Name] != "" {
		flagPtr.Name = "--" + parameterMap[Name]

		if p.PointersMap[flagPtr.Name] != nil {
			return errors.New("flag " + flagPtr.Name + " has been already defined")
		}

		p.PointersMap[flagPtr.Name] = flagPtr
	}

	p.Pointers = append(p.Pointers, flagPtr)
	return nil
}

func (p *Parser) ParseWithPipe(args []string, r io.Reader) error {
	builder := new(strings.Builder)

	for scanner := bufio.NewScanner(r); scanner.Scan(); {
		builder.WriteString(scanner.Text())
	}

	return p.Parse(append(args, builder.String()))
}

func (p *Parser) Parse(args []string) error {
	var current *FlagPointer

	req := p.buildUsage()

	for _, arg := range args {

		if (p.PointersMap["-h"] == nil && arg == "-h") || arg == "--help" {
			fmt.Println(p.Usage.String())
			os.Exit(0)
		}

		if current != nil {
			current.Parse(arg)

			req = slices.DeleteFunc(req, func(s *FlagPointer) bool {
				return s == current
			})

			current = nil
			continue
		} else if p.PointersMap[arg] == nil {
			fmt.Println(p.Usage.String())
			os.Exit(0)
		}

		if p.PointersMap[arg] != nil && !p.PointersMap[arg].IsBoolean() {
			current = p.PointersMap[arg]
		} else if p.PointersMap[arg] != nil {
			p.PointersMap[arg].Parse("true")
		}
	}

	if current != nil {
		fmt.Println(p.Usage.String())
		return errors.New("Unspecified value for " + current.FullName())
	}

	if len(req) > 0 {
		return errors.New("Flag " + req[0].FullName() + " is required")
	}

	return nil
}

func (p *Parser) buildUsage() []*FlagPointer {
	required := make([]*FlagPointer, 0)

	longestN := slices.MaxFunc(p.Pointers, func(a, b *FlagPointer) int {
		return len(a.N) - len(b.N)
	})

	longestName := slices.MaxFunc(p.Pointers, func(a, b *FlagPointer) int {
		return len(a.FullName()) - len(b.FullName())
	})

	for _, flag := range p.Pointers {
		if flag.Required {
			required = append(required, flag)
		}

		leftSpaces := LeftPadding + len(longestN.N)
		rightSpaces := RightPadding + len(longestName.FullName()) - len(flag.FullName())

		p.Usage.WriteString(fmt.Sprintln(""))
		p.Usage.WriteString(strings.Repeat(Space, 2))

		if flag.N == "" {
			rightSpaces -= leftSpaces
			p.Usage.WriteString(strings.Repeat(Space, leftSpaces))
		}

		p.Usage.WriteString(flag.FullName())
		p.Usage.WriteString(strings.Repeat(Space, rightSpaces))
		p.Usage.WriteString(flag.Usage)
	}

	return required
}
