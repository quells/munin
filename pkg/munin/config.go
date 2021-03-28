package munin

import (
	"bytes"
	"fmt"
	"math"
)

type GraphType int

const (
	Default GraphType = iota

	// Gauge is for things like temperature or number of people in a room.
	Gauge

	// Counter is for continuous incrementing counters which never decrease.
	// Overflows are taken into account by the Munin update function.
	Counter

	// Derive will store the derivative of the line going from the last to
	// the current value. This is like Counter but without overflow checks.
	Derive

	// Absolute is for counters which are reset upon reading.
	// For example, the number of messages received since the last poll.
	Absolute
)

func (t GraphType) String() string {
	switch t {
	case Gauge:
		return "GAUGE"
	case Counter:
		return "COUNTER"
	default:
		return ""
	}
}

// A Series is a single line on a Munin graph.
type Series struct {
	// Label is the human readable name for what the series represents.
	// It should be relatively short but descriptive.
	Label string

	// Info is a longer description of what the series represents.
	Info string

	// Type of value the series represents. See GraphType for descriptions.
	Type GraphType

	// Min value the series should have, for graph scaling.
	Min float64

	// Max value the series should have, for graph scaling.
	Max float64

	// Warn if the value for this series is above this value.
	Warn float64

	// Crit (critical) alarm if the value for this series is above this value.
	Crit float64
}

// NewSeries with a label and nothing else.
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

// Config values for a single Munin graph/plugin.
type Config struct {
	// Title of the graph.
	Title string

	// Category the graph will appear in.
	Category string

	// Info about what the graph is measuring.
	Info string

	// YAxis label (should include unit if applicable).
	YAxis string

	// Base for value scaling SI prefixes.
	// For example, if a value is a number of bytes then the base should be 1024
	// so that 1M represents 1048576 instead of 1000000, etc.
	Base int

	// Series which should be displayed on the graph, keyed by their internal field name.
	// These keys will be sanitized to meet Munin requirements.
	Series map[string]Series
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
		key = cleanFieldName(key)
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
