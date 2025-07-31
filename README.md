[![GoDoc](https://godoc.org/github.com/jessevdk/go-flags?status.png)](https://godoc.org/github.com/raiiga/flagg)
[![Go Report Card](https://goreportcard.com/badge/github.com/raiiga/flagg)](https://goreportcard.com/report/github.com/raiiga/flagg )
## Description:

Flagg is a CLI library for parsing command-line arguments in Go.

It is designed to be intuitive, easy to use and suitable for cases 
where usage of e.g. *spf13/cobra*, *urfave/cli*, *etc.* may be excessive.

## Example:

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
	opts, flags := new(options), flagg.New("Usage: app [options]")

	if err := flags.Map(opts); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*opts)
}
```
## Output:
```shell
  $ go build -o example ./main.go

$ ./example -h
Usage: app [options]
  -i, --input      Input file path
  -o, --output     Output file path
  -v, --version    Print version
  -V, --verbose    Verbose output mode
  
$ ./example --help
Usage: app [options]
  -i, --input      Input file path
  -o, --output     Output file path
  -v, --version    Print version
  -V, --verbose    Verbose output mode
  
$ ./example -i /etc/hosts -o out.log --verbose
{/etc/hosts out.log false true}
```
