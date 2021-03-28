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
