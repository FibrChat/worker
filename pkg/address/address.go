package address

import (
	"fmt"
	"strings"
)

// Address represents a user identity as id@domain.
type Address struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

func (a Address) String() string {
	return a.ID + "@" + a.Domain
}

// Parse parses a "user@domain" string into an Address.
func Parse(s string) (Address, error) {
	parts := strings.SplitN(s, "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Address{}, fmt.Errorf("invalid address %q: expected user@domain", s)
	}

	return Address{ID: parts[0], Domain: parts[1]}, nil
}
