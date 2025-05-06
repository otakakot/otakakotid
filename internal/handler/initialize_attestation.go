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

func InitializeAttestation(rw http.ResponseWriter, req *http.Request) {
	code, err := req.Cookie(core.CookeyRegistration)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	emailKV, err := cloudflare.NewKVNamespace(database.KVNSRegistrationEmail)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	email, err := emailKV.GetString(code.Value, nil)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	wa, err := core.NewWebAuthn(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	user := core.GenereteUser(email)

	creation, session, err := wa.BeginRegistration(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionKV, err := cloudflare.NewKVNamespace(database.KVNSWebAuthnSessionData)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	sid := uuid.NewString()

	sessionBuf := bytes.Buffer{}

	if err := json.NewEncoder(&sessionBuf).Encode(session); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := sessionKV.PutString(sid, base64.StdEncoding.EncodeToString(sessionBuf.Bytes()), nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	slog.InfoContext(req.Context(), "session base64: "+base64.StdEncoding.EncodeToString(sessionBuf.Bytes()))

	http.SetCookie(rw, &http.Cookie{
		Name:     core.CookeyAttestation,
		Value:    sid,
		Secure:   true,
		HttpOnly: true,
	})

	if err := json.NewEncoder(rw).Encode(creation.Response); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}
}
