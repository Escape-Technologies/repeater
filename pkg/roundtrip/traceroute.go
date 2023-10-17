package roundtrip

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/Escape-Technologies/repeater/pkg/logger"
	tr "github.com/aeden/traceroute"
)

func traceroute(input string) {
	logger.Debug("Traceroute debug info")
	u, err := url.Parse(input)
	if err != nil {
		logger.Debug("ERROR parsing url : %v", err)
		return
	}
	res, err := net.LookupHost(u.Hostname())
	if err != nil {
		logger.Debug("ERROR looking up host : %v", err)
		return
	}
	if len(res) == 0 {
		logger.Debug("No hosts found")
		return
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		port = 0
	}
	if port == 0 {
		if u.Scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	}

	logger.Debug("Starting traceroute to %v on port %v", res[0], port)
	trOpt := tr.TracerouteOptions{}
	trOpt.SetPort(port)
	trRes, err := tr.Traceroute(res[0], &trOpt)
	if err != nil {
		logger.Debug("ERROR running traceroute : %v", err)
		return
	}
	printTracerouteResult(trRes)
}

func printTracerouteResult(trRes tr.TracerouteResult) {
	logger.Debug(
		"IP : %v.%v.%v.%v",
		trRes.DestinationAddress[0],
		trRes.DestinationAddress[1],
		trRes.DestinationAddress[2],
		trRes.DestinationAddress[3],
	)
	for _, hop := range trRes.Hops {
		printHop(hop)
	}
}

func printHop(hop tr.TracerouteHop) {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	hostOrAddr := addr
	if hop.Host != "" {
		hostOrAddr = hop.Host
	}
	if hop.Success {
		logger.Debug("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime)
	} else {
		logger.Debug("%-3d *\n", hop.TTL)
	}
}
