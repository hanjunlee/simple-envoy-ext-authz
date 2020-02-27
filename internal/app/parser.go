package app

import (
	"fmt"
	"strings"
)

func parseToken(value string) (string, error) {
	const (
		bearer = "Bearer "
	)

	ss := strings.Split(value, bearer)
	if len(ss) != 2 {
		return "", fmt.Errorf("failed to parse: %s", value)
	}

	return ss[1], nil
}
