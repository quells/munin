package main

import (
	"strings"

	"github.com/quells/munin/internal/pihole5"
	"github.com/quells/munin/internal/set"
	"github.com/quells/munin/pkg/munin"
)

func main() {
	p := new(piHole)
	munin.Run(p)
}

type piHole struct{}

func (p *piHole) Help() string {
	return help
}

func (p *piHole) Config(env munin.Env) (conf munin.Config, err error) {
	conf.Title = "PiHole stats - " + env["host"]
	conf.Category = "dns"
	conf.Info = info
	conf.Series = make(map[string]munin.Series)

	set := skipSet(env)
	for k, label := range labels {
		if _, skip := set[k]; skip {
			continue
		}

		series := munin.NewSeries(label).
			WithType(munin.Gauge)
		if i, ok := infos[k]; ok {
			series = series.WithInfo(i)
		}
		conf.Series[k] = series
	}

	return
}

func (p *piHole) Fetch(env munin.Env) (values munin.Values, precision munin.Precision, err error) {
	client := pihole5.NewClient(env["host"], skipSet(env))
	values, precision, err = client.Load()
	return
}

func skipSet(env munin.Env) set.Strings {
	except := strings.Split(env["except"], ",")
	except = append(except, "ads_percentage_today")

	return set.OfStrings(except)
}
