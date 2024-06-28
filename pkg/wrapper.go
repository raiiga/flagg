package flagg

import (
	"flag"
	"os"
)

type DefaultWrapper struct {
	flags *flag.FlagSet
}

func Default() *DefaultWrapper {
	return &DefaultWrapper{
		flags: flag.NewFlagSet("default", flag.ContinueOnError),
	}
}

func (w *DefaultWrapper) NewMapper() *Mapper {
	return &Mapper{FlagSet: w.flags}
}

func (w *DefaultWrapper) NewHandler() *Handler {
	return &Handler{FlagMap: make(map[string]bool), FlagSet: w.flags}
}

func (w *DefaultWrapper) Parse() error {
	return w.flags.Parse(os.Args[1:])
}

type CustomWrapper struct {
	Flags *flag.FlagSet
}

func Custom(name string, handling flag.ErrorHandling) *CustomWrapper {
	return &CustomWrapper{flag.NewFlagSet(name, handling)}
}

func (w *CustomWrapper) NewMapper() *Mapper {
	return &Mapper{FlagSet: w.Flags}
}

func (w *CustomWrapper) NewHandler() *Handler {
	return &Handler{FlagMap: make(map[string]bool), FlagSet: w.Flags}
}

func (w *CustomWrapper) Parse(arguments []string) error {
	return w.Flags.Parse(arguments)
}
