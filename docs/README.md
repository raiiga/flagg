## Description

Flagg is simple wrapper for Go's flag package that provide full and short flags support, mapping to structure, or 
handle flags with defined functions.

## Usage

Mapping flags to defined structure:

```go
package main

import (
    "github.com/raiiga/flagg"
    "log"
)

type options struct {
    input   string `short:"i" long:"input" usage:"Input file path"`
    output  string `short:"o" long:"output" usage:"Output file path"`
    help    bool   `short:"h" long:"help" usage:"Show this message and exit"`
    version bool   `short:"v" long:"version" usage:"Show app version and exit"`
    verbose bool   `long:"verbose" usage:"Verbose output"`
}

func main() {
    mapper, opts := flagg.NewMapper("yourProjectName", flagg.ContinueOnError), options{}
    
    if err := mapper.Map(&opts); err != nil {
       log.Fatal(err)
    }
    
    if opts.help {
       mapper.FlagSet.Usage()
    } else {
       /* your further logic */
    }
}
```

Handling flags with defined functions:

```go
package main

import (
    "github.com/raiiga/flagg"
    "log"
)

func main() {
    handler := flagg.NewHandler("yourProjectName", flagg.ContinueOnError)
    
    handler.Func("input", "i", "Input file path", func(input string) error {
       /* your further logic */
       return nil
    })
    
    handler.Func("output", "o", "Output file path", func(input string) error {
       /* your further logic */
       return nil
    })

    if err := handler.Handle(); err != nil {
        log.Fatal(err)
    }
}
```