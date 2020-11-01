package util

import (
	"strings"
)

// ReplacerFromMap creates a strings.Replacer object from a map of orig => replacement
func ReplacerFromMap(args map[string]string) *strings.Replacer {
	// Convert into array of 0 => key, 1 => val, 2 => key, etc...
	arglist := make([]string, len(args)*2)
	c := 0
	for k, v := range args {
		arglist[(c * 2)] = k
		arglist[(c*2)+1] = v
		c++
	}
	return strings.NewReplacer(arglist...)
}
