package main

import (
	"net/http"

	"github.com/syumai/workers"

	"github.com/otakakot/otakakotid/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.Health)

	http.HandleFunc("/.well-known/openid-configuration", handler.WellKnown)

	workers.Serve(nil) // use http.DefaultServeMux
}
