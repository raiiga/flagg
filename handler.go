package flagg

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type handler struct {
	FlagMap map[string]bool
	FlagSet *flag.FlagSet
}

func NewHandler(name string, handling flag.ErrorHandling) *handler {
	return &handler{FlagMap: make(map[string]bool), FlagSet: flag.NewFlagSet(name, handling)}
}

func (h *handler) Func(fullKey, shortKey, usage string, f func(s string) error) {
	executor := h.checkedExecutor(fullKey, shortKey, f)

	if fullKey != emptyString {
		h.FlagSet.Func(fullKey, usage, executor)
	}

	if shortKey != emptyString {
		h.FlagSet.Func(shortKey, usage, executor)
	}
}

func (h *handler) BoolFunc(fullKey, shortKey, usage string, f func(s string) error) {
	executor := h.checkedExecutor(fullKey, shortKey, f)

	if fullKey != emptyString {
		h.FlagSet.BoolFunc(fullKey, usage, executor)
	}

	if shortKey != emptyString {
		h.FlagSet.BoolFunc(shortKey, usage, executor)
	}
}

func (h *handler) Handle() error {
	return h.HandleFromArgs(os.Args[1:])
}

func (h *handler) HandleFromArgs(args []string) error {
	return h.FlagSet.Parse(args)
}

func (h *handler) checkedExecutor(fullKey, shortKey string, f func(s string) error) func(s string) error {
	key := fmt.Sprintf(flagKeyFormat, fullKey, shortKey)

	return func(s string) error {
		if h.FlagMap[key] {
			return errors.New("only one flag may be used at a time")
		}

		h.FlagMap[key] = true

		return f(s)
	}
}
