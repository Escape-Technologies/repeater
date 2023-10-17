package roundtrip

import (
	"net"
	"net/url"

	"github.com/Escape-Technologies/repeater/pkg/logger"
	tr "github.com/pixelbender/go-traceroute/traceroute"
)

func traceroute(input string) {
	logger.Debug("Traceroute debug info")
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
		logger.Debug("No hosts found")
		return
	}
	hops, err := tr.Trace(net.ParseIP(res[0]))
	if err != nil {
		logger.Error("ERROR tracing ip %v", err)
	}
	for _, h := range hops {
		for _, n := range h.Nodes {
			logger.Debug("%d. %v %v", h.Distance, n.IP, n.RTT)
		}
	}
}
