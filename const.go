package flagg

import "flag"

const (
	emptyString   = ""
	flagKeyFormat = "%s:%s"

	longTagName  = "long"
	shortTagName = "short"
	valueTagName = "value"
	usageTagName = "usage"

	ContinueOnError = flag.ContinueOnError
	ExitOnError     = flag.ExitOnError
	PanicOnError    = flag.PanicOnError
)
