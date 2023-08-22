package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/elazarl/goproxy"
)

var UUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var url = "0.0.0.0:8080"

func main() {
	username := os.Getenv("ESCAPE_ORGANIZATION_ID")
	if !UUID.MatchString(username) {
		log.Printf("ESCAPE_ORGANIZATION_ID must be a UUID in lowercase")
		log.Printf("To get your organization id, go to https://app.escape.tech/organization/settings/")
		os.Exit(1)
	}
	password := os.Getenv("ESCAPE_REPEATER_ID")
	if !UUID.MatchString(password) {
		log.Printf("ESCAPE_REPEATER_ID must be a UUID in lowercase")
		log.Printf("To get your API key, go to https://app.escape.tech/user/profile/")
		os.Exit(1)
	}

	start(username, password)
}

func start(user string, pass string) {
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})

	total := url
	log.Printf("Listening on %s", total)
	log.Fatal(http.ListenAndServe(total, proxy))
}
