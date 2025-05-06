package handler

import (
	"encoding/json"
	"net/http"

	"github.com/otakakot/otakakotid/pkg/api"
)

func WellKnown(rw http.ResponseWriter, req *http.Request) {
	issuer := req.URL.Scheme + "://" + req.Host

	conf := api.OpenIDConfigurationResponseSchema{
		Issuer:                           issuer,
		AuthorizationEndpoint:            issuer + "/authorize",
		JwksUri:                          issuer + "/certs",
		RevocationEndpoint:               issuer + "/revoke",
		TokenEndpoint:                    issuer + "/token",
		UserinfoEndpoint:                 issuer + "/userinfo",
		SubjectTypesSupported:            []string{"public"},
		IdTokenSigningAlgValuesSupported: []string{"RS256"},
	}

	if err := json.NewEncoder(rw).Encode(conf); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
}
