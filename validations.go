package scalr

import (
	"net"
	"regexp"
	"strings"
)

// A regular expression used to validate common string ID patterns.
var reStringID = regexp.MustCompile(`^[a-zA-Z0-9\-\._]+$`)

// validString checks if the given input is present and non-empty.
func validString(v *string) bool {
	return v != nil && strings.TrimSpace(*v) != ""
}

// validStringID checks if the given string pointer is non-nil and
// contains a typical string identifier.
func validStringID(v *string) bool {
	return v != nil && reStringID.MatchString(*v)
}

func validIPv4Network(v *string) bool {
	if v == nil {
		return false
	}

	if ip := net.ParseIP(*v); ip != nil && ip.To4() != nil {
		return true
	}

	addr, _, err := net.ParseCIDR(*v)
	if err != nil {
		return false
	}

	return addr.To4() != nil
}
