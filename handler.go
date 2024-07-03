package flagg

import (
	"errors"
	"flag"
	"fmt"
)

const emptyString, flagKeyFormat = "", "%s:%s"

type handler struct {
	FlagMap map[string]bool
	FlagSet *flag.FlagSet
}

func (h *handler) Func(fullKey, shortKey string, f func(s string) error) {
	executor := h.checkedExecutor(fullKey, shortKey, f)

	if fullKey != emptyString {
		h.FlagSet.Func(fullKey, emptyString, executor)
	}

	if shortKey != emptyString {
		h.FlagSet.Func(shortKey, emptyString, executor)
	}
}

func (h *handler) BoolFunc(fullKey, shortKey string, f func(s string) error) {
	executor := h.checkedExecutor(fullKey, shortKey, f)

	if fullKey != emptyString {
		h.FlagSet.BoolFunc(fullKey, emptyString, executor)
	}

	if shortKey != emptyString {
		h.FlagSet.BoolFunc(shortKey, emptyString, executor)
	}
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
