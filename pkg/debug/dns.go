package debug

import (
	"net"
	"net/url"
	"strings"
)

func dns(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return err.Error()
	}
	res, err := net.LookupHost(u.Hostname())
	if err != nil {
		return err.Error()
	}
	if len(res) == 0 {
		return "No hosts found"
	}
	return strings.Join(res, "-")
}
