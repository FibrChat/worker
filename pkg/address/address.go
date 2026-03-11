package address

import (
	"fmt"
	"strings"
)

// Split splits an "user@domain" address into its components
func Split(addr string) (string, string, error) {
	parts := strings.SplitN(addr, "@", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", "", fmt.Errorf("invalid address %q: missing user", addr)
	}

	if parts[1] == "" {
		return "", "", fmt.Errorf("invalid address %q: missing domain", addr)
	}

	return parts[0], parts[1], nil
}
