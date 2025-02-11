## Description

Flagg is simple wrapper for Go's flag package that provide full and short flags support, mapping to structure, or 
handle flags with defined functions.

## Usage

Mapping flags to defined structure:

```go
package main

import (
	"fmt"
	"github.com/raiiga/flagg"
)

type options struct {
	input   string `flagg:"n:i, name:input, usage:Input file path"`
	output  string `flagg:"n:o, name:output, usage:Output file path"`
	version bool   `flagg:"n:v, name:version, usage:Print version"`
	verbose bool   `flagg:"n:V, name:verbose, usage:Verbose output mode"`
}

func main() {
	opts, flags := new(options), flagg.New("yourProjectName")

	if help, err := flags.Map(opts); err != nil {
		fmt.Println(err)
		return
	} else if help {
		return
	}

	fmt.Println(opts)
}
```

Handling flags with defined functions:

```go
package main

import (
	"fmt"
	"github.com/raiiga/flagg"
)

type options struct {
	boolHandler  func() error             `flagg:"n:b, name:bool, usage:function activated if flag provided"`
	plainHandler func(input string) error `flagg:"n:f, name:func, usage:function activated if flag with value provided"`
}

func main() {
	opts, flags := new(options), flagg.New("yourProjectName")

	opts.boolHandler = func() error {
		_, err := fmt.Println("Your boolean function output")
		return err
	}

	opts.plainHandler = func(input string) error {
		_, err := fmt.Printf("Your input is: %s%s", input, fmt.Sprintln())
		return err
	}

	if help, err := flags.Map(opts); err != nil {
		fmt.Println(err)
		return
	} else if help {
		return
	}
}
```