package env

import "strings"

func Parse(environ []string) (env map[string]string) {
	env = make(map[string]string)

	for _, e := range environ {
		parts := strings.Split(e, "=")
		key := parts[0]
		value := strings.Join(parts[1:], "=")
		env[key] = value
	}

	return
}
