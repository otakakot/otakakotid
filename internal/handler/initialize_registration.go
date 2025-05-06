package handler

import (
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/otakakot/otakakotid/internal/database"
	"github.com/otakakot/otakakotid/pkg/api"
	"github.com/syumai/workers/cloudflare"
)

func InitializeRegistration(rw http.ResponseWriter, req *http.Request) {
	body := api.InitialiseRegistrationJSONRequestBody{}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	kv, err := cloudflare.NewKVNamespace(database.KVNSRegistrationEmail)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	code := rand.Text()

	if err := kv.PutString(code, body.Email, nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	slog.InfoContext(req.Context(), "http://localhost:7777/registration?code="+code)

	// TODO: send email
}
