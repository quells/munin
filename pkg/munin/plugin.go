package munin

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/quells/munin/internal/env"
)

type Env map[string]string
type Values map[string]float64
type Precision map[string]int

type Plugin interface {
	Help() string
	Config(env Env) (conf Config, err error)
	Run(env Env) (values Values, precision Precision, err error)
}

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

func formatValue(value float64, precision int) string {
	if precision == 0 {
		return strconv.Itoa(int(math.Round(value)))
	}
	return strconv.FormatFloat(value, 'f', precision, 64)
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

func emitConfig(p Plugin, e Env) {
	conf, err := p.Config(e)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", conf)
}

func emitValues(p Plugin, e Env) {
	values, precision, err := p.Run(e)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	for k, v := range values {
		p := precision[k]
		buf.WriteString(k)
		buf.WriteString(".value ")
		buf.WriteString(formatValue(v, p))
		buf.WriteByte('\n')
	}
	fmt.Fprint(os.Stdout, buf.String())
}
