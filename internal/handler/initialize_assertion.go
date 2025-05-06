package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/otakakot/otakakotid/internal/core"
	"github.com/otakakot/otakakotid/internal/database"
	"github.com/syumai/workers/cloudflare"
)

func InitializeAssertion(rw http.ResponseWriter, req *http.Request) {
	wa, err := core.NewWebAuthn(req)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	assertion, session, err := wa.BeginDiscoverableLogin()
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sid := uuid.NewString()

	sessionBuf := bytes.Buffer{}

	if err := json.NewEncoder(&sessionBuf).Encode(session); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionKV, err := cloudflare.NewKVNamespace(database.KVNSWebAuthnSessionData)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := sessionKV.PutString(sid, base64.StdEncoding.EncodeToString(sessionBuf.Bytes()), nil); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     core.CookeyAssertion,
		Value:    sid,
		Secure:   true,
		HttpOnly: true,
	})

	if err := json.NewEncoder(rw).Encode(assertion.Response); err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}
}
