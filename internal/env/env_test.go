package env

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		environ []string
		want    map[string]string
	}{
		{
			"empty",
			[]string{},
			make(map[string]string),
		},
		{
			"some",
			[]string{"abc=123", "foo=bar=baz"},
			map[string]string{"abc": "123", "foo": "bar=baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.environ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
