package main

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/syumai/workers"

	"github.com/otakakot/otakakotid/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.Health)

	mux.HandleFunc("GET /assertion", handler.InitializeAssertion)
	mux.HandleFunc("POST /assertion", handler.FinalizeAssertion)

	mux.HandleFunc("POST /registration", handler.InitializeRegistration)
	mux.HandleFunc("GET /registration", handler.FinalizeRegistration)

	mux.HandleFunc("GET /attestation", handler.InitializeAttestation)
	mux.HandleFunc("POST /attestation", handler.FinalizeAttestation)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		Debug:            true,
	}).Handler(mux)

	workers.Serve(handler)
}
