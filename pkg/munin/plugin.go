// Copyright 2021 Kai Wells
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package munin

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/quells/munin/internal/env"
)

// Env variables passed in to the plugin.
type Env map[string]string

// Values produced by the plugin.
// Keyed by the field name of the corresponding Series.
type Values map[string]float64

// Precision (number of digits after the decimal place) for values produced by the plugin.
// Keyed by the field name of the corresponding Series.
type Precision map[string]int

// A Plugin for Munin which fits into this framework.
// The main function for a plugin should just have to pass in a value to the Run function
// and the rest will be taken care of.
type Plugin interface {
	// Help returns configuration and usage information for the plugin.
	// This does not conform to the Munin spec, but is useful since Go binaries
	// are not as easy to inspect as plaintext scripts.
	Help() string

	// Config returns graph configuration data, possibly informed by the environment
	// or configuration values passed in via the environment.
	Config(env Env) (conf Config, err error)

	// Fetch data values to be displayed on the graph, possibly informed by the environment
	// or configuration values passed in via the environment.
	Fetch(env Env) (values Values, precision Precision, err error)
}

// Run the Plugin as a good Munin citizen.
// Supports the "dirty config" capability for one-shot configuration and value emission.
func Run(p Plugin) {
	if helpRequested() {
		help := p.Help()
		fmt.Fprintf(os.Stdout, "%s\n", help)
		os.Exit(0)
	}

	e := env.Parse(os.Environ())

	if len(os.Args) == 2 && os.Args[1] == "config" {
		emitConfig(p, e)
		if e["MUNIN_CAP_DIRTYCONFIG"] == "1" {
			emitValues(p, e)
		}
		os.Exit(0)
	}

	emitValues(p, e)
	os.Exit(0)
}

func helpRequested() bool {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "help":
			return true
		case "--help":
			return true
		case "-h":
			return true
		}
	}

	return false
}

func formatValue(value float64, precision int) string {
	if precision == 0 {
		return strconv.Itoa(int(math.Round(value)))
	}
	return strconv.FormatFloat(value, 'f', precision, 64)
}

var fieldName = regexp.MustCompile(`(^[^A-Za-z_]|[^A-Za-z0-9_])`)

func cleanFieldName(text string) string {
	if text == "root" {
		return "_root"
	}

	return fieldName.ReplaceAllString(text, "_")
}

func emitConfig(p Plugin, e Env) {
	conf, err := p.Config(e)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", conf)
}

func emitValues(p Plugin, e Env) {
	values, precision, err := p.Fetch(e)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	for k, v := range values {
		p := precision[k]
		buf.WriteString(cleanFieldName(k))
		buf.WriteString(".value ")
		buf.WriteString(formatValue(v, p))
		buf.WriteByte('\n')
	}
	fmt.Fprint(os.Stdout, buf.String())
}
