package flagg

import (
	"flag"
	"os"
)

type Wrapper struct {
	flags *flag.FlagSet
}

// Default - creates and holds golang flag set value inside default *Wrapper instance.
// Initiates program exit if during flag processing occurs an error.
func Default() *Wrapper {
	return &Wrapper{flags: flag.NewFlagSet("default", flag.ExitOnError)}
}

// Custom - creates and holds golang flag set value inside custom *Wrapper instance.
// Initiates program exit if during flag processing occurs an error.
func Custom(name string, handling flag.ErrorHandling) *Wrapper {
	return &Wrapper{flag.NewFlagSet(name, handling)}
}

// NewMapper - creates new flag mapper instance.
// *mapper allows to parse and assign cli flags to given object value.
func (w *Wrapper) NewMapper() *mapper {
	return &mapper{FlagSet: w.flags}
}

// NewHandler - creates new flag handler instance.
// *handler allows to execute procedures when selected flag were used.
func (w *Wrapper) NewHandler() *handler {
	return &handler{FlagMap: make(map[string]bool), FlagSet: w.flags}
}

// Parse - initiates flag parsing process from standard (os.Args[1:]) arguments.
func (w *Wrapper) Parse() error {
	return w.flags.Parse(os.Args[1:])
}

// ParseFromArgs - initiates flag parsing process from given arguments.
func (w *Wrapper) ParseFromArgs(arguments []string) error {
	return w.flags.Parse(arguments)
}
