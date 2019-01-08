package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	targetHost   = os.Getenv("TARGET_HOST")
	targetSchema = os.Getenv("TARGET_SCHEMA")
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := &url.URL{
			Host:   targetHost,
			Scheme: targetSchema,
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = new(roundTripper)
		proxy.ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type roundTripper struct{}

func (roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("proxy to %q...", r.URL.String())
	return http.DefaultTransport.RoundTrip(r)
}
