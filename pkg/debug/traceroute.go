package debug

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	tr "github.com/aeden/traceroute"
)

func traceroute(input string) string {
	u, err := url.Parse(input)
	if err != nil {
		return "url.Parse : " + err.Error()
	}
	res, err := net.LookupHost(u.Hostname())
	if err != nil {
		return "net.LookupHost : " + err.Error()
	}
	if len(res) == 0 {
		return "No hosts found"
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

	trOpt := tr.TracerouteOptions{}
	trOpt.SetPort(port)
	trRes, err := tr.Traceroute(res[0], &trOpt)
	if err != nil {
		return "tr.Traceroute : " + err.Error()
	}
	return printTracerouteResult(trRes)
}

func printTracerouteResult(trRes tr.TracerouteResult) string {
	result := fmt.Sprintf(
		"%v.%v.%v.%v",
		trRes.DestinationAddress[0],
		trRes.DestinationAddress[1],
		trRes.DestinationAddress[2],
		trRes.DestinationAddress[3],
	)
	for _, hop := range trRes.Hops {
		result += printHop(hop)
	}
	return result
}

func printHop(hop tr.TracerouteHop) string {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	hostOrAddr := addr
	if hop.Host != "" {
		hostOrAddr = hop.Host
	}
	if hop.Success {
		return fmt.Sprintf("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime)
	}
	return fmt.Sprintf("%-3d *\n", hop.TTL)
}
