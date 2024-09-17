package stream

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/status"
)

func invalidClientID(msg string) bool { return msg == "invalid client UUID" }
func noSuchHost(msg string) bool      { return strings.Contains(msg, "no such host") }
func timedOut(msg string) bool        { return strings.Contains(msg, "i/o timeout") }

func grpcErrorFmt(s *status.Status) []string {
	res := []string{}
	if s == nil {
		return res
	}

	msg := s.Message()
	if invalidClientID(msg) {
		res = append(res, "Server rejected your configured client UUID (configuration issue)")
		res = append(res, "Check that the ESCAPE_REPEATER_ID environment variable is set correctly")
		res = append(res, "Go to https://app.escape.tech/organization/network/ to retrieve your repeaters list")

		return res
	}

	if noSuchHost(msg) {
		res = append(res, "Server could not be resolved (DNS issue)")
		res = append(res, "If you run inside docker, try passing the --network=host flag")
		res = append(res, "From the host machine, check that the server is reachable with nslookup repeater.escape.tech")
	}

	if timedOut(msg) {
		res = append(res, "Timed out connecting to the server (network issue)")
		res = append(res, "It may be linked to a firewall missconfiguration or a network issue")
		res = append(res, "Check that the documentation about firewall configuration https://docs.escape.tech/platform/enterprise/repeater#configure-your-firewall")
		res = append(res, "You can also try to run tracepath or tcptraceroute to check where the connection is blocked")
	}

	res = append(res, "If you need more help to debug this issue, please contact the support team with the result of a curl -v https://repeater.escape.tech")
	return res
}

func extractWhyStreamCreateError(err error) []string {
	res := []string{}
	if err == nil {
		return res
	}

	res = append(res, fmt.Sprintf("Error creating stream: %v", err))
	if s, ok := status.FromError(err); ok {
		res = append(res, grpcErrorFmt(s)...)
	}

	return res
}

func extractWhyRecvError(err error) []string {
	if err == nil {
		return []string{}
	}
	res := []string{}

	res = append(res, fmt.Sprintf("Error receiving data: %v", err))
	if s, ok := status.FromError(err); ok {
		res = append(res, grpcErrorFmt(s)...)
	}

	return res
}
