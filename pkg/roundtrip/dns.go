package roundtrip

import (
	"net"
	"net/url"

	"github.com/Escape-Technologies/repeater/pkg/logger"
)

func dns(input string) {
	logger.Debug("DNS debug info")
	u, err := url.Parse(input)
	if err != nil {
		logger.Error("ERROR parsing url : %v", err)
		return
	}
	res, err := net.LookupHost(u.Hostname())
	if err != nil {
		logger.Error("ERROR looking up host : %v", err)
		return
	}
	if len(res) == 0 {
		logger.Error("No hosts found")
		return
	}
	for _, r := range res {
		logger.Debug("Found dns result : %v", r)
	}
	return
}
