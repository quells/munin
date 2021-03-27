package munin

import (
	"bytes"
	"fmt"
	"math"
)

type GraphType int

const (
	Default GraphType = iota
	Gauge
)

func (t GraphType) String() string {
	switch t {
	case Gauge:
		return "GAUGE"
	default:
		return ""
	}
}

type Series struct {
	Label string
	Info  string
	Type  GraphType
	Min   float64
	Max   float64
	Warn  float64
	Crit  float64
}

func NewSeries(label string) (s Series) {
	s.Label = label
	s.Min = math.NaN()
	s.Max = math.NaN()
	s.Warn = math.NaN()
	s.Crit = math.NaN()
	return
}

func (s Series) WithInfo(info string) Series {
	s.Info = info
	return s
}

func (s Series) WithType(t GraphType) Series {
	s.Type = t
	return s
}

func (s Series) WithRange(min, max float64) Series {
	s.Min = min
	s.Max = max
	return s
}

func (s Series) WithWarnings(warn, crit float64) Series {
	s.Warn = warn
	s.Crit = crit
	return s
}

type Config struct {
	Title    string
	Category string
	Info     string
	YAxis    string
	Base     int
	Series   map[string]Series
}

func (c Config) String() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "graph_title %s\n", c.Title)
	if c.Base != 0 {
		fmt.Fprintf(buf, "graph_args --base %d\n", c.Base)
	}
	if c.Category != "" {
		fmt.Fprintf(buf, "graph_category %s\n", c.Category)
	}
	if c.YAxis != "" {
		fmt.Fprintf(buf, "graph_vlabel %s\n", c.YAxis)
	}
	if c.Info != "" {
		fmt.Fprintf(buf, "graph_info %s\n", c.Info)
	}

	for key, series := range c.Series {
		if series.Label != "" {
			fmt.Fprintf(buf, "%s.label %s\n", key, series.Label)
		}
		if series.Type != Default {
			fmt.Fprintf(buf, "%s.type %s\n", key, series.Type)
		}
		if !math.IsNaN(series.Min) {
			fmt.Fprintf(buf, "%s.min %f\n", key, series.Min)
		}
		if !math.IsNaN(series.Max) {
			fmt.Fprintf(buf, "%s.max %f\n", key, series.Max)
		}
		if !math.IsNaN(series.Warn) {
			fmt.Fprintf(buf, "%s.warn %f\n", key, series.Warn)
		}
		if !math.IsNaN(series.Crit) {
			fmt.Fprintf(buf, "%s.crit %f\n", key, series.Crit)
		}
		if series.Info != "" {
			fmt.Fprintf(buf, "%s.info %s\n", key, series.Info)
		}
	}

	return buf.String()
}
