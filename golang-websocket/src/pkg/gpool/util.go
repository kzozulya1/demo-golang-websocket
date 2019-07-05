package gpool

import (
	"strings"
)

// HTTPisClosedConnError reports whether err is an error from use of a closed
// network connection.
// Code is taken from https://golang.org/src/net/http/h2_bundle.go?h=use+of+closed+network+connection
func HTTPisClosedConnError(err error) bool {
	if err == nil {
		return false
	}

	str := err.Error()
	if strings.Contains(str, "use of closed network connection") {
		return true
	}
	return false
}
