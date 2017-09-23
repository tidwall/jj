package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

var (
	version = "0.0.1"
	tag     = "jsoned - JSON Stream Editor " + version
	usage   = `
usage: jsoned [-v value] [-r] [-D] [-O] [-p] [-i infile] [-o outfile] keypath

examples: jsoned keypath                      read value from stdin
      or: jsoned -i infile keypath            read value from infile
      or: jsoned -v value keypath             edit value
      or: jsoned -v value -o outfile keypath  edit value and write to outfile

options:
      -v value             Edit JSON key path value
      -r                   Use raw values, otherwise types are auto-detected
      -O                   Performance boost for value updates.
	  -D                   Delete the value at the specified key path
	  -p                   Make json pretty
      -i infile            Use input file instead of stdin
      -o outfile           Use output file instead of stdout
      keypath              JSON key path (like "name.last")

for more info: https://github.com/tidwall/jsoned
`
)

type args struct {
	infile    *string
	outfile   *string
	value     *string
	raw       bool
	del       bool
	opt       bool
	keypathok bool
	keypath   string
	pretty    bool
}

func parseArgs() args {
	fail := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, "%s\n", tag)
		if format != "" {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		}
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		os.Exit(1)
	}
	help := func() {
		buf := &bytes.Buffer{}
		fmt.Fprintf(os.Stderr, "%s\n", tag)
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		os.Stdout.Write(buf.Bytes())
		os.Exit(0)
	}
	var a args
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		default:
			if !a.keypathok {
				a.keypathok = true
				a.keypath = os.Args[i]
			} else {
				fail("unknown option argument: \"%s\"", a.keypath)
			}
		case "-v", "-i", "-o":
			arg := os.Args[i]
			i++
			if i >= len(os.Args) {
				fail("argument missing after: \"%s\"", arg)
			}
			switch arg {
			case "-v":
				a.value = &os.Args[i]
			case "-i":
				a.infile = &os.Args[i]
			case "-o":
				a.outfile = &os.Args[i]
			}
		case "-r":
			a.raw = true
		case "-p":
			a.pretty = true
		case "-D":
			a.del = true
		case "-O":
			a.opt = true
		case "-h", "--help", "-?":
			help()
		}
	}
	if !a.keypathok && !a.pretty {
		fail("missing required option: \"keypath\"")
	}
	return a
}

func main() {
	a := parseArgs()
	var input []byte
	var err error
	var outb []byte
	var outs string
	var f *os.File
	if a.infile == nil {
		input, err = ioutil.ReadAll(os.Stdin)
	} else {
		input, err = ioutil.ReadFile(*a.infile)
	}
	if err != nil {
		goto fail
	}
	if a.del {
		outb, err = sjson.DeleteBytes(input, a.keypath)
		if err != nil {
			goto fail
		}
	} else if a.value != nil {
		raw := a.raw
		val := *a.value
		if !raw {
			switch val {
			default:
				if len(val) > 0 {
					if (val[0] >= '0' && val[0] <= '9') || val[0] == '-' {
						if _, err := strconv.ParseFloat(val, 64); err == nil {
							raw = true
						}
					}
				}
			case "true", "false", "null":
				raw = true
			}
		}
		opts := &sjson.Options{}
		if a.opt {
			opts.Optimistic = true
			opts.ReplaceInPlace = true
		}
		if raw {
			// set as raw block
			outb, err = sjson.SetRawBytesOptions(
				input, a.keypath, []byte(val), opts)
		} else {
			// set as a string
			outb, err = sjson.SetBytesOptions(input, a.keypath, val, opts)
		}
		if err != nil {
			goto fail
		}
	} else {
		if !a.keypathok {
			outb = input
		} else {
			res := gjson.GetBytes(input, a.keypath)
			if a.raw {
				outs = res.Raw
			} else {
				outs = res.String()
			}
		}
	}
	if a.outfile == nil {
		f = os.Stdout
	} else {
		f, err = os.Create(*a.outfile)
		if err != nil {
			goto fail
		}
	}

	if outb != nil {
		if a.pretty {
			f.Write(pretty.Pretty(outb))
		} else {
			f.Write(outb)
		}
	} else {
		if a.pretty {
			f.Write(pretty.Pretty([]byte(outs)))
		} else {
			f.WriteString(outs)
		}
	}
	f.Close()
	return
fail:
	fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
	os.Exit(1)
}
