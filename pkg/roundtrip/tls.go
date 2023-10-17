package roundtrip

import (
	"net/url"

	cryptotls "crypto/tls"

	"github.com/Escape-Technologies/repeater/pkg/logger"
)

func tls(input string) {
	logger.Debug("TLS debug info")
	u, err := url.Parse(input)
	if err != nil {
		logger.Error("ERROR parsing url : %v", err)
		return
	}
	conf := &cryptotls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := cryptotls.Dial("tcp", u.Host, conf)
	if err != nil {
		logger.Error("Error in Dial %v", err)
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		logger.Debug("Issuer Name: %s", cert.Issuer)
		logger.Debug("Expiry: %s", cert.NotAfter.Format("2006-January-02"))
		logger.Debug("Common Name: %s", cert.Issuer.CommonName)
	}
}
