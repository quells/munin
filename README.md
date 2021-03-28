# Munin Plugin Framework for Go

A collection of helpers for writing [Munin](http://munin-monitoring.org/) plugins with some examples.

## Examples

[Pi-Hole stats](https://github.com/quells/munin/tree/main/cmd/pihole)

## Basic Usage

[Random number](https://github.com/quells/munin/tree/main/cmd/pihole/example/example.go)

```go
package main

import (
	"math/rand"

	"github.com/quells/munin/pkg/munin"
)

func main() {
	munin.Run(new(myPlugin))
}

type myPlugin struct{}

func (p *myPlugin) Help() string {
	return "My Munin Plugin"
}

func (p *myPlugin) Config(env munin.Env) (conf munin.Config, err error) {
	conf.Title = "My Data"
	conf.Info = "This is just an example."
	conf.Series = make(map[string]munin.Series)

	conf.Series["example"] = munin.NewSeries("data").WithType(munin.Gauge)

	return
}

func (p *myPlugin) Fetch(env munin.Env) (values munin.Values, precision munin.Precision, err error) {
	values = make(munin.Values)
	precision = make(munin.Precision)

	// return a random number between 0.000 and 0.999
	values["example"] = rand.Float64()
	precision["example"] = 3

	return
}
```
